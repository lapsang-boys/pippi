package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/decomp/exp/bin"
	_ "github.com/decomp/exp/bin/elf" // register ELF decoder
	"github.com/kr/pretty"
	"github.com/pkg/errors"
)

func usage() {
	const use = `
Parse ELF files.

Usage:

	alvriket [OPTION]... FILE

Flags:
`
	fmt.Fprintln(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	// Parse command line arguments.
	var (
		// Path to binary executable.
		binPath string
	)
	flag.Usage = usage
	flag.Parse()
	switch flag.NArg() {
	case 0:
		// Read from standard input.
		binPath = "-"
	case 1:
		binPath = flag.Arg(0)
	default:
		flag.Usage()
		os.Exit(1)
	}

	if err := parseELF(binPath); err != nil {
		log.Fatalf("%+v", err)
	}
}

// parseELF parses the given ELF file.
func parseELF(binPath string) error {
	file, err := bin.ParseFile(binPath)
	if err != nil {
		return errors.WithStack(err)
	}
	pretty.Println(file)
	return nil
}
