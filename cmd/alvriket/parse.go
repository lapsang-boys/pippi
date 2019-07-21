package main

import (
	"context"
	"flag"

	"github.com/decomp/exp/bin"
	"github.com/google/subcommands"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
)

// parseCmd is the command to parse an ELF file from the command line.
type parseCmd struct{}

func (*parseCmd) Name() string {
	return "parse"
}

func (*parseCmd) Synopsis() string {
	return "parse ELF file from command line"
}

func (*parseCmd) Usage() string {
	const use = `
Parse ELF file from command line.

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
	elfPath := f.Arg(0)

	// Connect to gRPC server.
	if err := parse(elfPath); err != nil {
		warn.Printf("parse failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// parse parses the given ELF file, pretty-printing its contents to standard
// output.
func parse(elfPath string) error {
	// Parse ELF file.
	file, err := bin.ParseFile(elfPath)
	if err != nil {
		return errors.WithStack(err)
	}
	// Pretty-print to standard output.
	pretty.Println(file)
	return nil
}
