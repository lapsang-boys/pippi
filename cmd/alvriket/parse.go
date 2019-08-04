package main

import (
	"context"
	"flag"

	"github.com/decomp/exp/bin"
	"github.com/google/subcommands"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
)

// parseCmd is the command to parse a binary file from the command line.
type parseCmd struct{}

func (*parseCmd) Name() string {
	return "parse"
}

func (*parseCmd) Synopsis() string {
	return "parse binary file from command line"
}

func (*parseCmd) Usage() string {
	const use = `
Parse binary file from command line.

Usage:
	parse [OPTION]... FILE

Flags:
`
	return use[1:]
}

func (cmd *parseCmd) SetFlags(f *flag.FlagSet) {
}

func (cmd *parseCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Parse command arguments.
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	binPath := f.Arg(0)
	// Parse binary file and pretty-print to standard output.
	if err := parse(binPath); err != nil {
		warn.Printf("parse failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// parse parses the given binary file, pretty-printing its contents to standard
// output.
func parse(binPath string) error {
	// Parse binary file.
	file, err := bin.ParseFile(binPath)
	if err != nil {
		return errors.WithStack(err)
	}
	// Pretty-print to standard output.
	pretty.Println(file)
	return nil
}
