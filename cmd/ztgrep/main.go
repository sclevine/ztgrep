package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"

	"github.com/sclevine/ztgrep"
)

type Options struct {
	Search struct {
		SkipBody bool `short:"b" long:"skip-body" description:"Skip file bodies"`
		SkipName bool `short:"n" long:"skip-name" description:"Skip file names inside of tarballs"`
		MaxZipSize int64 `short:"z" long:"max-zip-size" default:"0" default-mask:"10 MB" description:"Maximum zip file size to search"`
	} `group:"Search Options"`

	General struct {
		Version bool `short:"v" long:"version" description:"Return ztgrep version"`
	} `group:"General Options"`
}

var (
	Version = "0.0.0"
	opts Options
)

func main() {
	log.SetFlags(0)

	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassAfterNonOption|flags.PassDoubleDash)
	parser.Usage = "[OPTIONS] regexp paths..."
	restArgs, err := parser.Parse()
	if err != nil {
		if err, ok := err.(*flags.Error); ok && err.Type == flags.ErrHelp {
			log.Fatal(err)
		}
		log.Fatalf("Invalid arguments: %s", err)
	}
	if opts.General.Version {
		fmt.Printf("ztgrep v%s\n", Version)
		os.Exit(0)
	}
	if len(restArgs) == 0 {
		parser.WriteHelp(os.Stderr)
		os.Exit(0)
	}
	if len(restArgs) == 1 {
		restArgs = append(restArgs, "-")
	}
	if err := grep(restArgs[0], restArgs[1:]); err != nil {
		log.Fatalf("Failed: %s", err)
	}
}

func grep(expr string, paths []string) error {
	zt, err := ztgrep.New(expr)
	if err != nil {
		return err
	}
	if opts.Search.MaxZipSize != 0 {
		zt.MaxZipSize = opts.Search.MaxZipSize
	}
	zt.SkipName = opts.Search.SkipName
	zt.SkipBody = opts.Search.SkipBody
	for res := range zt.Start(paths) {
		path := strings.Join(res.Path, ":")
		if res.Err != nil {
			log.Printf("ztgrep: %s: %s", path, res.Err)
		} else {
			fmt.Println(path)
		}
	}
	return nil
}
