package tzgrep

import (
	"archive/tar"
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

func New(expr string) (*TZgrep, error) {
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
	Out                chan Result
	SkipName, SkipBody bool
	exp                *regexp.Regexp
}

type Result struct {
	Path []string
	Err  error
}

func (tz *TZgrep) Start(paths []string) {
	// TODO: restrict number of open files
	// TODO: buffer output to guarantee order
	wg := sync.WaitGroup{}
	wg.Add(len(paths))
	go func() {
		wg.Wait()
		close(tz.Out)
	}()
	for _, p := range paths {
		p := p
		go func() {
			tz.findPath(p)
			wg.Done()
		}()
	}
}

func (tz *TZgrep) findPath(path string) {
	if path == "-" {
		tz.find(os.Stdin, []string{"-"})
		return
	}
	f, err := os.Open(path)
	if err != nil {
		tz.Out <- Result{Path: []string{path}, Err: err}
	}
	defer f.Close()
	tz.find(f, []string{path})
}

func (tz *TZgrep) find(zr io.Reader, path []string) {
	zf, isTar := newDecompressor(path[len(path)-1])
	if !isTar && tz.SkipBody {
		return
	}
	r, err := zf(zr)
	if err != nil {
		tz.Out <- Result{Path: path, Err: err}
	}
	defer r.Close()

	if !isTar {
		if tz.exp.MatchReader(bufio.NewReader(r)) {
			tz.Out <- Result{Path: path}
		}
		return
	}

	tr := tar.NewReader(r)
	for h, err := tr.Next(); err != io.EOF; h, err = tr.Next() {
		if err != nil {
			tz.Out <- Result{Path: path, Err: err}
			break
		}
		path = append(path[:len(path):len(path)], h.Name)
		if !tz.SkipName {
			if tz.exp.MatchString(h.Name) {
				tz.Out <- Result{Path: path}
			}
		}
		tz.find(tr, path)
	}
}

type decompressor func(io.Reader) (io.ReadCloser, error)

func newDecompressor(path string) (zf decompressor, isTar bool) {
	p := strings.ToLower(path)
	switch {
	case hasSuffixes(p, ".tar.gz", ".tgz", ".taz"):
		return gzReader, true
	case hasSuffixes(p, ".tar.bz2", ".tar.bz", ".tbz", ".tbz2", ".tz2", ".tb2"):
		return bz2Reader, true
	case hasSuffixes(p, ".tar.xz", ".txz"):
		return xzReader, true
	case hasSuffixes(p, ".tar.zst", ".tzst", ".tar.zstd"):
		return zstReader, true
	case hasSuffixes(p, ".tar"):
		return nopReader, true

	case hasSuffixes(p, ".gz"):
		return gzReader, false
	case hasSuffixes(p, ".bz2", ".bz"):
		return bz2Reader, false
	case hasSuffixes(p, ".xz"):
		return xzReader, false
	case hasSuffixes(p, ".zst", ".zstd"):
		return zstReader, false
	default:
		return nopReader, false
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
