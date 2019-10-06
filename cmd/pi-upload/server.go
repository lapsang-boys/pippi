package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net"
	"os"

	"github.com/google/subcommands"
	"github.com/lapsang-boys/pippi/pkg/pi"
	uploadpb "github.com/lapsang-boys/pippi/proto/upload"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	// Default gRPC address to listen on.
	grpcAddr = ":1100"
)

// serverCmd is the command to launch a gRPC server processing parse binary
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
Reply to binary file parse request.

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

// Maximum file size in bytes.
const maxFileSize = 1 * 1024 * 1024 * 1024 // 1 GB

// listen listens on the given gRPC address for incoming requests to parse binary
// files.
func listen(addr string) error {
	dbg.Printf("listening on %q", addr)
	// Launch gRPC server.
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.WithStack(err)
	}
	server := grpc.NewServer(grpc.MaxRecvMsgSize(maxFileSize))
	// Register binary parser service.
	uploadpb.RegisterUploadServer(server, &uploadServer{})
	if err := server.Serve(l); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// uploadServer handles binary file upload requests.
type uploadServer struct{}

// Upload handles binary file upload requests.
func (us uploadServer) Upload(ctx context.Context, req *uploadpb.UploadRequest) (*uploadpb.UploadReply, error) {
	dbg.Printf("receiving %q", req.Filename)
	// Compute binary ID based on binary file contents.
	binID := pi.BinID(req.Content)
	if binID != req.Hash {
		return nil, errors.Errorf("hash mismatch; expected %q, got %q", req.Hash, binID)
	}
	// Create project directory for the binary ID.
	binDir, err := pi.BinDir(binID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return nil, errors.WithStack(err)
	}
	// Store binary file to disk.
	binPath, err := pi.BinPath(binID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := ioutil.WriteFile(binPath, req.Content, 0644); err != nil {
		return nil, errors.WithStack(err)
	}
	// Send reply.
	reply := &uploadpb.UploadReply{
		Id: binID,
	}
	return reply, nil
}
