package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jessevdk/go-flags"

	"github.com/sclevine/tzgrep"
)

type Options struct {
	SkipBody bool `short:"b" long:"skip-body" description:"Skip file bodies"`
	SkipName bool `short:"n" long:"skip-name" description:"Skip file names inside of tarballs"`
}

var opts Options

func main() {
	log.SetFlags(0)

	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassAfterNonOption|flags.PassDoubleDash)
	restArgs, err := parser.Parse()
	if err != nil {
		if err, ok := err.(*flags.Error); ok && err.Type == flags.ErrHelp {
			log.Fatal(err)
		}
		log.Fatalf("Invalid arguments: %s", err)
	}
	if len(restArgs) == 0 {
		log.Fatal("Missing expression")
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
	tz.SkipName = opts.SkipName
	tz.SkipBody = opts.SkipBody
	tz.Start(paths)
	for res := range tz.Out {
		if res.Err != nil {
			log.Printf("tzgrep: %s", res.Err)
		} else {
			fmt.Println(strings.Join(res.Path, ":"))
		}
	}
	return nil
}
