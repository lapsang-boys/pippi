package main

import (
	"context"
	"flag"
	"net"

	"github.com/decomp/exp/bin"
	"github.com/google/subcommands"
	"github.com/lapsang-boys/pippi/alvriket"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	// Default gRPC address to listen on.
	grpcAddr = ":1234"
)

// serverCmd is the command to launch a gRPC server processing parse ELF
// requests.
type serverCmd struct {
	// gRPC address to listen on.
	Addr string
}

func (*serverCmd) Name() string {
	return "server"
}

func (*serverCmd) Synopsis() string {
	return "launch gRPC server"
}

func (*serverCmd) Usage() string {
	const use = `
Reply to ELF file parse request.

Usage:
	server [OPTION]...

Flags:
`
	return use[1:]
}

func (cmd *serverCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.Addr, "addr", grpcAddr, "gRPC address to listen on")
}

func (cmd *serverCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := listen(cmd.Addr); err != nil {
		warn.Printf("listen failed; %+v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// listen listens on the given gRPC address for incoming requests to parse ELF
// files.
func listen(addr string) error {
	dbg.Printf("listening on %q", addr)
	// Launch gRPC server.
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.WithStack(err)
	}
	server := grpc.NewServer()
	// Register ELF parser service.
	alvriket.RegisterELFParserServer(server, &elfParserServer{})
	if err := server.Serve(l); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// elfParserServer implements alvriket.ELFParserServer.
type elfParserServer struct{}

// ParseELF parses the given ELF file.
func (elfParserServer) ParseELF(ctx context.Context, req *alvriket.ParseELFRequest) (*alvriket.ParseELFReply, error) {
	dbg.Printf("parsing %q", req.ElfPath)
	// Parse ELF file.
	file, err := bin.ParseFile(req.ElfPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// Send reply.
	reply := &alvriket.ParseELFReply{
		Nsects: int32(len(file.Sections)),
	}
	return reply, nil
}
