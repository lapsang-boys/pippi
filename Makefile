all: pippigo pippirust pippipython

pippigo:
	go install ./...
	make -C cmd/pippi
	go mod tidy

pippirust:
	make -C cmd/strings

pippipython:
	# nothing to do. TODO: type-check using Cython?

run_frontend:
	make -C cmd/pippi run

run_backend:
	forego start

clean:
	make -C cmd/pippi clean
	make -C cmd/strings clean

.PHONY: all pippigo pippirust pippipython
