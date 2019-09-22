all: pippigo pippirust pippipython

pippigo:
	go install ./...

pippirust:
	make -C cmd/strings

pippipython:
	# nothing to do. TODO: type-check using Cython?

clean:
	make -C cmd/strings clean

.PHONY: all pippigo pippirust pippipython
