PROTO=$(wildcard *.proto)
DIR=$(basename ${PROTO})
PB_GO=$(PROTO:.proto=.pb.go)
PB_PY=$(PROTO:.proto=_pb2.py) $(PROTO:.proto=_pb2_grpc.py)
PB_RS=$(PROTO:.proto=.rs) $(PROTO:.proto=_grpc.rs)

all: ${DIR} ${PB_GO} ${PB_PY} ${PB_RS}

%: %.proto
	mkdir -p $@

# Protobuf compiler:
#
#    $ pacman -S protoc

# Go protobuf compiler:
#
#    $ go get -u github.com/golang/protobuf/protoc-gen-go

%.pb.go: %.proto
	protoc -I . --go_out=plugins=grpc:$(basename $<) $<

# Python protobuf compiler:
#
#    $ python -m pip install --user grpcio
#    $ python -m pip install --user grpcio-tools

%_pb2.py %_pb2_grpc.py: %.proto
	python -m grpc_tools.protoc -I . --python_out=$(basename $<) --grpc_python_out=$(basename $<) $<

# Rust protobuf compiler:
#
#    $ cargo install protobuf-codegen
#    $ cargo install grpcio-compiler

%.rs %_grpc.rs: %.proto
	protoc -I . --rust_out=$(basename $<) --grpc_out=$(basename $<) --plugin=protoc-gen-grpc=`which grpc_rust_plugin` $<

clean:
	$(RM) -v ${PB_GO} ${PB_PY} ${PB_RS}
