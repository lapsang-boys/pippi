# Sockerdricka is a tool for disassembling binary executables.
#
# The disassembler is based on the "Shingled Graph Disassembly" method [1],
# which also gave inspiration to the name. One way to think about the Shingled
# Graph Disassembly method is that it tries to create disassembly trees rooted
# at every byte-offset of the binary executable, and then invalidate those
# starting byte-offsets for trees in which leaf nodes correspond to invalid
# assembly encodedings. Pruning those paths, what is left is but a tree filled
# with Sockerdricka.
#
# [1]: https://link.springer.com/chapter/10.1007/978-3-319-06608-0_23

# Dependencies:
#
#    $ python -m pip install capstone
#    $ python -m pip install click
#    $ python -m pip install hexdump

# gRPC
import grpc
import sys
sys.path.append('../../proto/bin') # Include pippi import paths.
import bin_pb2 as bin
import bin_pb2_grpc as bin_grpc

# Command line handling.
import click

# Get cache directory.
import appdirs
import pathlib

# Hex dump.
import hexdump

#from capstone import *

# Binary file extension.
ext = ".bin"

@click.command()
@click.option('--bin_server', default='127.0.0.1:1234', help='Binary parser server.', show_default=True)
@click.argument('bin_id', type=str)
def client_cmd(bin_server, bin_id):
	# Validate binary ID
	if not valid_id(bin_id):
		print('invalid ID; expected lowercase sha256 hash, got %s' % (bin_id))
		return
	# Get binary file path.
	cache_dir = appdirs.user_cache_dir('pippi', 'lapsang-boys')
	bin_name = bin_id + ext
	bin_path = pathlib.Path(cache_dir).joinpath(bin_id, bin_name)
	bin_data = bin_path.read_bytes()

	with grpc.insecure_channel(bin_server) as channel:
		# Parse binary file.
		client = bin_grpc.BinaryParserStub(channel)
		file = client.ParseBinary(bin.ParseBinaryRequest(bin_id=bin_id))
		# Locate executable sections.
		for sect in file.sections:
			if bin.X in sect.perms:
				print(sect)
				hexdump.hexdump(bin_data[sect.offset:sect.offset+sect.length])

# Hexadecimal digits.
hex = "0123456789abcdef"

# valid_id reports whether the given binary ID is valid (lowercase sha256 hash).
def valid_id(bin_id):
	if len(bin_id) != 64:
		return False
	if bin_id.lower() != bin_id:
		return False
	for c in bin_id:
		if c not in hex:
			return False
	return True

if __name__ == '__main__':
	client_cmd()

#cs = Cs(CS_ARCH_X86, CS_MODE_32)
#for i in cs.disasm(data, 0x1000):
#	print("0x%x:\t%s\t%s" % (i.address, i.mnemonic, i.op_str))
