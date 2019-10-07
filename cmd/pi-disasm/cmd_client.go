package main

import (
	"context"
	"flag"

	"github.com/decomp/exp/bin"
	"github.com/google/subcommands"
	"github.com/kr/pretty"
	"github.com/lapsang-boys/pippi/pkg/pi"
	binpb "github.com/lapsang-boys/pippi/proto/bin"
	disasmpb "github.com/lapsang-boys/pippi/proto/disasm"
	"github.com/mewkiz/pkg/jsonutil"
	"github.com/mewkiz/pkg/pathutil"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// clientCmd is the command to connect to a gRPC server to send disassemble
// binary requests.
type clientCmd struct {
	// path instruction addresses JSON file.
	instAddrJSONPath string
	// disasm gRPC address to connect to.
	disasmAddr string
}

func (*clientCmd) Name() string {
	return "client"
}

func (*clientCmd) Synopsis() string {
	return "connect to gRPC server"
}

func (*clientCmd) Usage() string {
	const use = `
Send disassemble binary file request.

Usage:
	client [OPTION]... BIN_ID

Flags:
`
	return use[1:]
}

func (cmd *clientCmd) SetFlags(f *flag.FlagSet) {
	// TODO: add -arch flag.
	f.StringVar(&cmd.instAddrJSONPath, "inst_addrs", "", "path instruction addresses JSON file")
	f.StringVar(&cmd.disasmAddr, "addr_disasm", defaultDisasmAddr, "disasm gRPC address to connect to")
}

func (cmd *clientCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Parse command arguments.
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	binID := f.Arg(0)
	// Connect to gRPC server.
	if err := connect(cmd.instAddrJSONPath, cmd.disasmAddr, binID); err != nil {
		warn.Printf("connect failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// connect connects to the given gRPC address to send a disassemble binary file
// reuquest.
func connect(instAddrJSONPath, disasmAddr, binID string) error {
	if err := pi.CheckBinID(binID); err != nil {
		return errors.WithStack(err)
	}
	dbg.Printf("connecting to %q", disasmAddr)
	binPath, err := pi.BinPath(binID)
	if err != nil {
		return errors.WithStack(err)
	}
	// If -inst_addrs flag was not specified, look for instruction addresses JSON
	// file in "{BIN_DIR}/{BIN_ID}.json".
	if len(instAddrJSONPath) == 0 {
		instAddrJSONPath = pathutil.TrimExt(binPath) + ".json"
	}
	var instAddrs []bin.Address
	dbg.Printf("parsing %q", instAddrJSONPath)
	if err := jsonutil.ParseFile(instAddrJSONPath, &instAddrs); err != nil {
		return errors.WithStack(err)
	}
	// Connect to gRPC server.
	conn, err := grpc.Dial(disasmAddr, grpc.WithInsecure())
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()
	// Send disassemble binary file request.
	client := disasmpb.NewDisassemblerClient(conn)
	ctx := context.Background()
	var instAddrspb []uint64
	for _, instAddr := range instAddrs {
		instAddrspb = append(instAddrspb, uint64(instAddr))
	}
	req := &disasmpb.DisassembleRequest{
		BinId:     binID,
		InstAddrs: instAddrspb,
		Arch:      binpb.Arch_X86_64, // TODO: make configurable.
	}
	reply, err := client.Disassemble(ctx, req)
	if err != nil {
		return errors.WithStack(err)
	}
	pretty.Println(reply)
	return nil
}
