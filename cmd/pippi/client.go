package main

import (
	"context"
	"log"
	"time"

	binpb "github.com/lapsang-boys/pippi/proto/bin"
	disasmpb "github.com/lapsang-boys/pippi/proto/disasm"
	disasm_objdumppb "github.com/lapsang-boys/pippi/proto/disasm_objdump"
	stringspb "github.com/lapsang-boys/pippi/proto/strings"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func Strings(stringsAddr, binId string) ([]*stringspb.StringInfo, error) {
	// Connect to gRPC server.
	conn, err := grpc.Dial(stringsAddr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()

	// Send binary parsing request.
	client := stringspb.NewStringsExtractorClient(conn)
	ctx := context.Background()

	req := &stringspb.StringsRequest{
		BinId: binId,
	}
	now := time.Now()
	reply, err := client.ExtractStrings(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println(time.Since(now))

	return reply.Strings, nil
}

func Sections(binAddr, binId string) (*binpb.File, error) {
	// Connect to gRPC server.
	conn, err := grpc.Dial(binAddr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()

	// Send binary parsing request.
	client := binpb.NewBinaryParserClient(conn)
	ctx := context.Background()

	req := &binpb.ParseBinaryRequest{
		BinId: binId,
	}
	now := time.Now()
	file, err := client.ParseBinary(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println(time.Since(now))

	return file, nil
}

func Disassembly(disasmAddr, binId string, arch binpb.Arch, instAddrs []uint64) (*disasmpb.DisassembleReply, error) {
	// Connect to gRPC server.
	conn, err := grpc.Dial(disasmAddr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()

	// Send binary parsing request.
	client := disasmpb.NewDisassemblerClient(conn)
	ctx := context.Background()

	req := &disasmpb.DisassembleRequest{
		BinId:     binId,
		Arch:      arch,
		InstAddrs: instAddrs,
	}
	now := time.Now()
	reply, err := client.Disassemble(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println(time.Since(now))

	return reply, nil
}

func InstAddrs(disasmObjdumpAddr, binId string) ([]uint64, error) {
	// Connect to gRPC server.
	conn, err := grpc.Dial(disasmObjdumpAddr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()

	// Send binary parsing request.
	client := disasm_objdumppb.NewInstAddrExtractorClient(conn)
	ctx := context.Background()

	req := &disasm_objdumppb.InstAddrsRequest{
		BinId: binId,
	}
	now := time.Now()
	reply, err := client.ExtractInstAddrs(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println(time.Since(now))

	return reply.InstAddrs, nil
}
