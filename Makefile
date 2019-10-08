all: protogen pippigo pippirust pippipython testdatagen

protogen:
	make -C proto

pippigo:
	go install ./...
	make -C cmd/pippi
	go mod tidy

pippirust:
	# TODO: uncomment when #43 is resolved.
	#make -C cmd/pi-strings

pippipython:
	# nothing to do. TODO: type-check using Cython?

testdatagen:
	make -C testdata

run_frontend:
	make -C cmd/pippi run

run_backend:
	forego start

clean:
	make -C cmd/pippi clean
	make -C cmd/pi-strings clean
	make -C proto clean
	make -C testdata clean

.PHONY: all protogen pippigo pippirust pippipython
