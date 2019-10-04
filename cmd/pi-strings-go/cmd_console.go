package main

import (
	"context"
	"flag"
	"io/ioutil"

	"github.com/google/subcommands"
	"github.com/kr/pretty"
)

// consoleCmd is the command to parse a binary file from the command line.
type consoleCmd struct{}

func (*consoleCmd) Name() string {
	return "console"
}

func (*consoleCmd) Synopsis() string {
	return "extract strings of binary files from command line"
}

func (*consoleCmd) Usage() string {
	const use = `
Extract strings of binary files from command line.

Usage:
	console [OPTION]... FILE

Flags:
`
	return use[1:]
}

func (cmd *consoleCmd) SetFlags(f *flag.FlagSet) {
}

func (cmd *consoleCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Parse command arguments.
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	binPath := f.Arg(0)
	// Extract strings of binary file and pretty-print to standard output.
	buf, err := ioutil.ReadFile(binPath)
	if err != nil {
		warn.Printf("unable to read contents of binary file %q; %+v", binPath, err)
		return subcommands.ExitFailure
	}
	// TODO: figure out a good way to make minimum length configurable.
	infos := extractStrings(buf, defaultMinLength)
	pretty.Println(infos)
	return subcommands.ExitSuccess
}
