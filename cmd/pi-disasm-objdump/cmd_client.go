package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/kr/pretty"
	"github.com/lapsang-boys/pippi/pkg/pi"
	disasm_objdumppb "github.com/lapsang-boys/pippi/proto/disasm_objdump"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// clientCmd is the command to connect to a gRPC server to send instruction
// address extraction requests.
type clientCmd struct {
	// disasm_objdump gRPC address to listen on.
	disasmObjdumpAddr string
}

func (*clientCmd) Name() string {
	return "client"
}

func (*clientCmd) Synopsis() string {
	return "connect to gRPC server"
}

func (*clientCmd) Usage() string {
	const use = `
Send to instruction address extraction requests.

Usage:
	client [OPTION]... BIN_ID

Flags:
`
	return use[1:]
}

func (cmd *clientCmd) SetFlags(f *flag.FlagSet) {
	// TODO: add -arch flag.
	f.StringVar(&cmd.disasmObjdumpAddr, "addr_disasm_objdump", defaultDisasmObjdumpAddr, "disasm_objdump gRPC address to connect to")
}

func (cmd *clientCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Parse command arguments.
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	binID := f.Arg(0)
	// Connect to gRPC server.
	if err := connect(cmd.disasmObjdumpAddr, binID); err != nil {
		warn.Printf("connect failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// connect connects to the given gRPC address to send an extract instruction
// addresses reuquest.
func connect(disasmObjdumpAddr, binID string) error {
	if err := pi.CheckBinID(binID); err != nil {
		return errors.WithStack(err)
	}
	dbg.Printf("connecting to %q", disasmObjdumpAddr)
	// Connect to gRPC server.
	conn, err := grpc.Dial(disasmObjdumpAddr, grpc.WithInsecure())
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()
	// Send extract instruction addresses request.
	client := disasm_objdumppb.NewInstAddrExtractorClient(conn)
	ctx := context.Background()
	req := &disasm_objdumppb.InstAddrsRequest{
		BinId: binID,
	}
	reply, err := client.ExtractInstAddrs(ctx, req)
	if err != nil {
		return errors.WithStack(err)
	}
	pretty.Println(reply)
	return nil
}
