package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sclevine/tzgrep"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 3 {
		log.Fatal("Invalid args")
	}
	if err := grep(os.Args[1], os.Args[2:]); err != nil {
		log.Fatalf("Failed: %s", err)
	}
}

func grep(expr string, paths []string) error {
	tz, err := tzgrep.New(expr)
	if err != nil {
		return err
	}
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
