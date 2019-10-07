package main

import (
	"context"
	"flag"
	"net"

	"github.com/google/subcommands"
	"github.com/lapsang-boys/pippi/pkg/pi"
	disasm_objdumppb "github.com/lapsang-boys/pippi/proto/disasm_objdump"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	// Default disasm gRPC address to listen on.
	disasmObjdumpGRPCAddr = ":1310"
)

// serverCmd is the command to launch a gRPC server processing instruction
// address extraction requests.
type serverCmd struct {
	// disasm_objdump gRPC address to listen on.
	DisasmObjdumpAddr string
}

func (*serverCmd) Name() string {
	return "server"
}

func (*serverCmd) Synopsis() string {
	return "launch gRPC server"
}

func (*serverCmd) Usage() string {
	const use = `
Reply to instruction address extraction requests.

Usage:
	server [OPTION]...

Flags:
`
	return use[1:]
}

func (cmd *serverCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.DisasmObjdumpAddr, "addr_disasm_objdump", disasmObjdumpGRPCAddr, "disasm_objdump gRPC address to listen on")
}

func (cmd *serverCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := listen(cmd.DisasmObjdumpAddr); err != nil {
		warn.Printf("listen failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// listen listens on the given gRPC address for incoming requests to extract
// instruction addresses.
func listen(disasmObjdumpAddr string) error {
	dbg.Printf("listening on %q", disasmObjdumpAddr)
	// Launch gRPC server.
	l, err := net.Listen("tcp", disasmObjdumpAddr)
	if err != nil {
		return errors.WithStack(err)
	}
	server := grpc.NewServer()
	// Register disassembler service.
	disasm_objdumppb.RegisterInstAddrExtractorServer(server, &disasmObjdumpServer{})
	if err := server.Serve(l); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// disasmObjdumpServer implements disasm_objdumppb.InstAddrExtractorServer.
type disasmObjdumpServer struct {
	// bin gRPC address.
	binAddr string
}

// ExtractInstAddr extract instruction address of the given binary file.
func (s *disasmObjdumpServer) ExtractInstAddrs(ctx context.Context, req *disasm_objdumppb.InstAddrsRequest) (*disasm_objdumppb.InstAddrsReply, error) {
	if err := pi.CheckBinID(req.BinId); err != nil {
		return nil, errors.WithStack(err)
	}
	dbg.Printf("extracting instruction addresses of ID %q", req.BinId)
	binPath, err := pi.BinPath(req.BinId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Extract instruction addresses.
	instAddrs, err := extractInstAddrs(binPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Send reply.
	reply := &disasm_objdumppb.InstAddrsReply{}
	for _, instAddr := range instAddrs {
		reply.InstAddrs = append(reply.InstAddrs, uint64(instAddr))
	}
	return reply, nil
}
