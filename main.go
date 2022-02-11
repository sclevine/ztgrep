package main

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 3 {
		log.Fatal("Invalid args")
	}
	if err := tzgrep(os.Args[1], os.Args[2:]); err != nil {
		log.Fatalf("Failed: %s", err)
	}
}

func tzgrep(expr string, paths []string) error {
	tz, err := NewTZgrep(expr)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(paths))
	go func() {
		wg.Wait()
		tz.Close()
	}()
	for _, p := range paths {
		p := p
		go func() {
			tz.FindPath(p)
			wg.Done()
		}()
	}
	for res := range tz.Out {
		if res.Err != nil {
			log.Printf("tzgrep: %s", res.Err)
		} else {
			fmt.Println(strings.Join(res.Path, ":"))
		}
	}
	return nil
}

func NewTZgrep(expr string) (*TZgrep, error) {
	exp, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &TZgrep{
		Out: make(chan Result),
		exp: exp,
	}, nil
}

type TZgrep struct {
	Out chan Result
	exp *regexp.Regexp
}

type Result struct {
	Path []string
	Err  error
}

func (tz *TZgrep) Close() {
	close(tz.Out)
}

func (tz *TZgrep) FindPath(path string) {
	f, err := os.Open(path)
	if err != nil {
		tz.Out <- Result{Path: []string{path}, Err: err}
	}
	defer f.Close()
	tz.Find(f, []string{path})
}

func (tz *TZgrep) Find(zr io.Reader, path []string) {
	if tz.exp.MatchString(path[len(path)-1]) {
		tz.Out <- Result{Path: path}
	}
	zf, isTar := newDecompressor(path[len(path)-1])
	if !isTar {
		return
	}
	r, err := zf(zr)
	if err != nil {
		tz.Out <- Result{Path: path, Err: err}
	}
	defer r.Close()
	tr := tar.NewReader(r)
	for h, err := tr.Next(); err != nil; h, err = tr.Next() {
		tz.Find(tr, append(path[:len(path):len(path)], h.Name))
	}
}

type decompressor func(io.Reader) (io.ReadCloser, error)

func newDecompressor(path string) (zf decompressor, ok bool) {
	p := strings.ToLower(path)
	switch {
	case hasSuffixes(p, ".tar"):
		return func(r io.Reader) (io.ReadCloser, error) {
			return io.NopCloser(r), nil
		}, true
	case hasSuffixes(p, ".tar.gz", ".tgz", ".taz"):
		return func(r io.Reader) (io.ReadCloser, error) {
			r, err := gzip.NewReader(r)
			return io.NopCloser(r), err
		}, true
	case hasSuffixes(p, ".tar.bz2", ".tar.bz", ".tbz", ".tbz2", ".tz2", ".tb2"):
		return func(r io.Reader) (io.ReadCloser, error) {
			return io.NopCloser(bzip2.NewReader(r)), nil
		}, true
	case hasSuffixes(p, ".tar.xz", ".txz"):
		return xzReader, true
	case hasSuffixes(p, ".tar.zst", ".tzst", ".tar.zstd"):
		return zstdReader, true
	default:
		return nil, false
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

func xzReader(r io.Reader) (io.ReadCloser, error) {
	return zCmdReader(exec.Command("xz", "-d", "-T0"), r)
}

func zstdReader(r io.Reader) (io.ReadCloser, error) {
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
