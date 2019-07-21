package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
	uploadpb "github.com/lapsang-boys/pippi/proto/upload"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	// Default gRPC address to listen on.
	grpcAddr  = ":1235"
	extension = ".bin"
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

// listen listens on the given gRPC address for incoming requests to parse binary
// files.
func listen(addr string) error {
	dbg.Printf("listening on %q", addr)
	// Launch gRPC server.
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.WithStack(err)
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return errors.WithStack(err)
	}
	cacheDir = filepath.Join(cacheDir, "pippi")

	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return errors.WithStack(err)
	}

	server := grpc.NewServer()
	// Register binary parser service.
	uploadpb.RegisterUploadServer(server, &uploadServer{cacheDir: cacheDir})
	if err := server.Serve(l); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type uploadServer struct {
	cacheDir string
}

func (us uploadServer) Upload(ctx context.Context, req *uploadpb.UploadRequest) (*uploadpb.UploadReply, error) {
	dbg.Printf("receiving %q", req.Filename)

	rawHash := sha256.Sum256(req.Content)
	hash := hex.EncodeToString(rawHash[:])
	if hash != req.Hash {
		return nil, errors.Errorf("Hash mismatch. Expected %s, got %s", req.Hash, hash)
	}
	dirName := filepath.Join(us.cacheDir, req.Hash)
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	destPath := filepath.Join(dirName, hash+extension)

	err = ioutil.WriteFile(destPath, req.Content, 0644)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Send reply.
	reply := &uploadpb.UploadReply{
		Id: hash,
	}
	return reply, nil
}
