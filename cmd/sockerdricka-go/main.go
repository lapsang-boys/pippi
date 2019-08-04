//go:generate protoc -I ../../proto --go_out=plugins=grpc:../../proto/disasm ../../proto/disasm.proto

package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/google/subcommands"
	"github.com/mewkiz/pkg/term"
)

var (
	// dbg is a logger with the "sockerdricka:" prefix which logs debug messages
	// to standard error.
	dbg = log.New(os.Stderr, term.CyanBold("sockerdricka:")+" ", 0)
	// warn is a logger with the "sockerdricka:" prefix which logs warning
	// messages to standard error.
	warn = log.New(os.Stderr, term.RedBold("sockerdricka:")+" ", 0)
)

func main() {
	// Initialize subcommands.
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&serverCmd{}, "")
	subcommands.Register(&clientCmd{}, "")
	subcommands.Register(&disasmCmd{}, "")

	// Parse command line arguments.
	flag.Parse()
	// Run subcommand based on command line arguments.
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
