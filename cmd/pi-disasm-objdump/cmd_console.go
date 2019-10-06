package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/mewkiz/pkg/jsonutil"
)

// consoleCmd is the command to extract instruction addresses of binary files
// from the command line.
type consoleCmd struct{}

func (*consoleCmd) Name() string {
	return "console"
}

func (*consoleCmd) Synopsis() string {
	return "extract instruction address of binary files from command line"
}

func (*consoleCmd) Usage() string {
	const use = `
Extract instruction addresses of binary files from command line.

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
	instAddrs, err := extractInstAddrs(binPath)
	if err != nil {
		warn.Printf("extraction instruction addresses from %q failed; %+v", binPath, err)
		return subcommands.ExitFailure
	}
	if err := jsonutil.Write(os.Stdout, instAddrs); err != nil {
		warn.Printf("unable to output instruction addresses; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
