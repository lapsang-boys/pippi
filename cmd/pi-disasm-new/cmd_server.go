package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/decomp/exp/bin"
	"github.com/google/subcommands"
	"github.com/lapsang-boys/pippi/pkg/pi"
	"github.com/lapsang-boys/pippi/pkg/services"
	binpb "github.com/lapsang-boys/pippi/proto/bin"
	disasmpb "github.com/lapsang-boys/pippi/proto/disasm"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	// Default bin gRPC address to listen on.
	defaultBinAddr = fmt.Sprintf("localhost:%d", services.BinPort)
	// Default disasm gRPC address to listen on.
	defaultDisasmAddr = fmt.Sprintf("localhost:%d", services.DisasmPort)
)

// serverCmd is the command to launch a gRPC server processing disassemble
// binary requests.
type serverCmd struct {
	// bin gRPC address to listen on.
	binAddr string
	// disasm gRPC address to listen on.
	disasmAddr string
}

func (*serverCmd) Name() string {
	return "server"
}

func (*serverCmd) Synopsis() string {
	return "launch gRPC server"
}

func (*serverCmd) Usage() string {
	const use = `
Reply to disassemble binary file request.

Usage:
	server [OPTION]...

Flags:
`
	return use[1:]
}

func (cmd *serverCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.binAddr, "addr_bin", defaultBinAddr, "bin gRPC address to listen on")
	f.StringVar(&cmd.disasmAddr, "addr_disasm", defaultDisasmAddr, "disasm gRPC address to connect to")
}

func (cmd *serverCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := listen(cmd.binAddr, cmd.disasmAddr); err != nil {
		warn.Printf("listen failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// listen listens on the given gRPC address for incoming requests to disassemble
// binary files.
func listen(binAddr, disasmAddr string) error {
	dbg.Printf("listening on %q", disasmAddr)
	// Launch gRPC server.
	l, err := net.Listen("tcp", disasmAddr)
	if err != nil {
		return errors.WithStack(err)
	}
	server := grpc.NewServer()
	// Register disassembler service.
	disasmpb.RegisterDisassemblerServer(server, &disasmServer{binAddr: binAddr})
	if err := server.Serve(l); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// disasmServer implements disasmpb.DisassembleServer.
type disasmServer struct {
	// bin gRPC address.
	binAddr string
}

// Disassemble disassembles the given binary file.
func (s *disasmServer) Disassemble(ctx context.Context, req *disasmpb.DisassembleRequest) (*disasmpb.DisassembleReply, error) {
	if err := pi.CheckBinID(req.BinId); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("disassembling ID %q", req.BinId)
	// Parse instruction addresses.
	db := &Database{}
	for _, instAddr := range req.InstAddrs {
		db.InstAddrs = append(db.InstAddrs, bin.Address(instAddr))
	}
	// Parse binary file.
	binPath, err := pi.BinPath(req.BinId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	insts, err := disasmBinary(db, binArch(req.Arch), binPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Send reply.
	reply := &disasmpb.DisassembleReply{}
	for _, inst := range insts {
		i := &disasmpb.Instruction{
			Addr:    uint64(inst.Addr()),
			InstStr: inst.String(),
		}
		reply.Insts = append(reply.Insts, i)
	}
	return reply, nil
}

// binArch returns the bin.Arch machine architecture corresponding to the given
// protobuf enum.
func binArch(arch binpb.Arch) bin.Arch {
	// TODO: move into pkg/pi/bin?
	switch arch {
	case binpb.Arch_X86_32:
		return bin.ArchX86_32
	case binpb.Arch_X86_64:
		return bin.ArchX86_64
	case binpb.Arch_MIPS_32:
		return bin.ArchMIPS_32
	case binpb.Arch_PowerPC_32:
		return bin.ArchPowerPC_32
	default:
		panic(fmt.Errorf("support for machine architecture %v not yet implemented", uint64(arch)))
	}
}
