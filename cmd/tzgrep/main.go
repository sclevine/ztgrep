package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"

	"github.com/sclevine/tzgrep"
)

type Options struct {
	Search struct {
		SkipBody bool `short:"b" long:"skip-body" description:"Skip file bodies"`
		SkipName bool `short:"n" long:"skip-name" description:"Skip file names inside of tarballs"`
	} `group:"Search Options"`
}

var opts Options

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
	tz, err := tzgrep.New(expr)
	if err != nil {
		return err
	}
	tz.SkipName = opts.Search.SkipName
	tz.SkipBody = opts.Search.SkipBody
	tz.Start(paths)
	for res := range tz.Out {
		path := strings.Join(res.Path, ":")
		if res.Err != nil {
			log.Printf("tzgrep: %s: %s", path, res.Err)
		} else {
			fmt.Println(path)
		}
	}
	return nil
}
