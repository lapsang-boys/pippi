package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/decomp/exp/bin"
	"github.com/google/subcommands"
	"github.com/mewkiz/pkg/jsonutil"
)

// consoleCmd is the command to parse a binary file from the command line.
type consoleCmd struct {
	// path instruction addresses JSON file.
	InstAddrJSONPath string
}

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
	console [OPTION]... -inst_addrs FILE.json FILE

Flags:
`
	return use[1:]
}

func (cmd *consoleCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.InstAddrJSONPath, "inst_addrs", "", "path instruction addresses JSON file")
}

func (cmd *consoleCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Parse command arguments.
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	binPath := f.Arg(0)
	if len(cmd.InstAddrJSONPath) == 0 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	db := &Database{}
	if err := jsonutil.ParseFile(cmd.InstAddrJSONPath, &db.InstAddrs); err != nil {
		warn.Printf("parse JSON file %q failed; %+v", cmd.InstAddrJSONPath, err)
		return subcommands.ExitFailure
	}
	arch := bin.ArchX86_64 // TODO: make configurable.
	insts, err := disasmBinary(db, arch, binPath)
	if err != nil {
		warn.Printf("unable to disassemble binary %q; %+v", binPath, err)
		return subcommands.ExitFailure
	}
	for _, inst := range insts {
		fmt.Printf("0x%08X    %s\n", uint64(inst.Addr()), inst.String())
	}
	return subcommands.ExitSuccess
}
