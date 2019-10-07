//+build ignore

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/google/subcommands"
	"github.com/kr/pretty"
	"github.com/lapsang-boys/pippi/cmd/pi-bin/binpbx"
	"github.com/lapsang-boys/pippi/pkg/pi"
	binpb "github.com/lapsang-boys/pippi/proto/bin"
	disasmpb "github.com/lapsang-boys/pippi/proto/disasm"
	"github.com/pkg/errors"
	"golang.org/x/arch/x86/x86asm"
	"google.golang.org/grpc"
)

// clientCmd is the command to connect to a gRPC server to send disassemble
// binary requests.
type clientCmd struct {
	// bin gRPC address to connect to.
	binAddr string
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
	f.StringVar(&cmd.binAddr, "addr_bin", defaultBinAddr, "bin gRPC address to connect to")
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
	if err := connect(cmd.binAddr, cmd.disasmAddr, binID); err != nil {
		warn.Printf("connect failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// connect connects to the given gRPC address to send a disassemble binary file
// reuquest.
func connect(binAddr, disasmAddr, binID string) error {
	if err := pi.CheckBinID(binID); err != nil {
		return errors.WithStack(err)
	}
	dbg.Printf("connecting to %q", disasmAddr)
	// Connect to gRPC server.
	conn, err := grpc.Dial(disasmAddr, grpc.WithInsecure())
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()
	// Parse binary file.
	file, err := binpbx.ParseFile(binAddr, binID)
	if err != nil {
		return errors.WithStack(err)
	}
	// Send disassemble binary file request.
	client := disasmpb.NewDisassemblerClient(conn)
	ctx := context.Background()
	req := &disasmpb.DisassembleRequest{
		BinId: binID,
	}
	reply, err := client.Disassemble(ctx, req)
	if err != nil {
		return errors.WithStack(err)
	}
	pretty.Println(reply)
	// Processor mode (16-, 32-, or 64-bit exection mode).
	var mode int
	switch file.Arch {
	case binpb.Arch_X86_32:
		mode = 32
	case binpb.Arch_X86_64:
		mode = 64
	}
	// TODO: use receive.proto to get file from server. Right now, we assume that
	// we are running on localhost to read the file contents of binID.
	// Read file contents.
	binPath, err := pi.BinPath(binID)
	if err != nil {
		return errors.WithStack(err)
	}
	binData, err := ioutil.ReadFile(binPath)
	if err != nil {
		return errors.WithStack(err)
	}
	// Print disassembly.
	for _, execSect := range reply.ExecSections {
		for _, off := range execSect.ValidOffsets {
			sectData := binData[execSect.Section.Offset : execSect.Section.Offset+execSect.Section.Length]
			instAddr := execSect.Section.Addr + off
			inst, err := x86asm.Decode(sectData[off:], mode)
			if err != nil {
				panic(fmt.Errorf("invalid instruction reported as valid at address 0x%08X", instAddr))
			}
			fmt.Printf("%08X\t%v\n", instAddr, inst)
		}
	}
	return nil
}
