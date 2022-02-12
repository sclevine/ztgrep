package ztgrep

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/sync/semaphore"
)

const defaultMaxZipSize = 10 << (10 * 2) // 10 MB

var cpuLock = semaphore.NewWeighted(int64(runtime.NumCPU()))

func acquireCPU() { cpuLock.Acquire(context.Background(), 1) }
func releaseCPU() { cpuLock.Release(1) }

// New returns a *ZTgrep given a regular expression following https://golang.org/s/re2syntax
func New(expr string) (*ZTgrep, error) {
	exp, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &ZTgrep{
		MaxZipSize: defaultMaxZipSize,
		exp:        exp,
	}, nil
}

// ZTgrep searchs for file names and contents within nested compressed archives.
type ZTgrep struct {
	MaxZipSize int64 // maximum size of zip file to search (held in memory)
	SkipName   bool  // skip file names
	SkipBody   bool  // skip file contents

	exp *regexp.Regexp
}

// Result contains each matching path in Path.
// Each entry in Path[1:] represents a file nested in the previous archive.
type Result struct {
	Path []string
	Err  error
}

// Start searches paths in parallel, returning results via a channel
func (zt *ZTgrep) Start(paths []string) <-chan Result {
	// TODO: restrict number of open files
	// TODO: buffer output to guarantee order
	out := make(chan Result)
	wg := sync.WaitGroup{}
	wg.Add(len(paths))
	go func() {
		wg.Wait()
		close(out)
	}()
	go func() {
		for _, p := range paths {
			p := p
			acquireCPU() // encourge ordered output
			go func() {
				zt.findPath(out, p)
				releaseCPU()
				wg.Done()
			}()
		}
	}()
	return out
}

func (zt *ZTgrep) findPath(out chan<- Result, path string) {
	if path == "-" {
		zt.find(out, os.Stdin, []string{"-"})
		return
	}
	f, err := os.Open(path)
	if err != nil {
		out <- Result{Path: []string{path}, Err: err}
		return
	}
	defer f.Close()
	zt.find(out, f, []string{path})
}

// TODO: implement version that uses file headers to identify type
func (zt *ZTgrep) find(out chan<- Result, zr io.Reader, path []string) {
	zf, xf := zt.newDecompressor(path[len(path)-1])
	if xf == nil && zt.SkipBody {
		return
	}
	r, err := zf(zr)
	if err != nil {
		out <- Result{Path: path, Err: err}
		return
	}
	defer r.Close()

	if xf == nil {
		if zt.exp.MatchReader(bufio.NewReader(r)) {
			out <- Result{Path: path}
		}
		return
	}

	if err := xf(r, func(name string, fr io.Reader) error {
		p := append(path[:len(path):len(path)], name)
		if !zt.SkipName {
			if zt.exp.MatchString(name) {
				out <- Result{Path: p}
			}
		}
		zt.find(out, fr, p)
		return nil
	}); err != nil {
		out <- Result{Path: path, Err: err}
		return
	}
}

type decompressor func(io.Reader) (io.ReadCloser, error)

type extractor func(io.Reader, func(string, io.Reader) error) error

func (zt *ZTgrep) newDecompressor(path string) (zf decompressor, xf extractor) {
	p := strings.ToLower(path)
	switch {
	case hasSuffixes(p, ".tar.gz", ".tgz", ".taz"):
		return gzReader, tarReader
	case hasSuffixes(p, ".tar.bz2", ".tar.bz", ".tbz", ".tbz2", ".tz2", ".tb2"):
		return bz2Reader, tarReader
	case hasSuffixes(p, ".tar.xz", ".txz"):
		return xzReader, tarReader
	case hasSuffixes(p, ".tar.zst", ".tzst", ".tar.zstd"):
		return zstReader, tarReader
	case hasSuffixes(p, ".tar"):
		return nopReader, tarReader
	case hasSuffixes(p, ".zip"):
		return nopReader, zt.zipReader

	case hasSuffixes(p, ".gz"):
		return gzReader, nil
	case hasSuffixes(p, ".bz2", ".bz"):
		return bz2Reader, nil
	case hasSuffixes(p, ".xz"):
		return xzReader, nil
	case hasSuffixes(p, ".zst", ".zstd"):
		return zstReader, nil
	default:
		return nopReader, nil
	}
}

func hasSuffixes(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

func tarReader(r io.Reader, fn func(string, io.Reader) error) error {
	tr := tar.NewReader(r)
	for h, err := tr.Next(); err != io.EOF; h, err = tr.Next() {
		if err != nil {
			return err
		}
		if err := fn(h.Name, tr); err != nil {
			return err
		}
	}
	return nil
}

func (zt *ZTgrep) zipReader(r io.Reader, fn func(string, io.Reader) error) error {
	tr, err := zt.readZip(r)
	if err != nil {
		return err
	}
	for _, file := range tr.File {
		fr, err := file.Open()
		if err != nil {
			return err // TODO: process next file if alg error?
		}
		if err := fn(file.Name, fr); err != nil {
			return err
		}
	}
	return nil
}

func nopReader(r io.Reader) (io.ReadCloser, error) {
	return io.NopCloser(r), nil
}

func gzReader(r io.Reader) (io.ReadCloser, error) {
	r, err := gzip.NewReader(r)
	return io.NopCloser(r), err
}

func bz2Reader(r io.Reader) (io.ReadCloser, error) {
	return io.NopCloser(bzip2.NewReader(r)), nil
}

func xzReader(r io.Reader) (io.ReadCloser, error) {
	return zCmdReader(exec.Command("xz", "-d", "-T0"), r)
}

func zstReader(r io.Reader) (io.ReadCloser, error) {
	return zCmdReader(exec.Command("zstd", "-d"), r)
}

func zCmdReader(cmd *exec.Cmd, r io.Reader) (io.ReadCloser, error) {
	cmd.Stdin = r
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return splitCloser{out, closerFunc(func() error {
		return cmd.Wait()
	})}, nil
}

type closerFunc func() error

func (f closerFunc) Close() error {
	return f()
}

type splitCloser struct {
	io.Reader
	io.Closer
}

func (zt *ZTgrep) readZip(r io.Reader) (*zip.Reader, error) {
	if f, ok := r.(*os.File); ok && f != os.Stdin {
		if fi, err := f.Stat(); err == nil {
			return zip.NewReader(f, fi.Size())
		}
	}
	limitedReader := &io.LimitedReader{R: r, N: zt.MaxZipSize}
	data, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}
	if limitedReader.N <= 0 {
		return nil, errors.New("zip file larger than limit")
	}
	br := bytes.NewReader(data)
	return zip.NewReader(br, br.Size())
}
