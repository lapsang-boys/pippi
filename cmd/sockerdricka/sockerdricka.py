# Sockerdricka is a tool for disassembling binary executables.
#
# The disassembler is based on the "Shingled Graph Disassembly" method [1],
# which also gave inspiration to the name. One way to think about the Shingled
# Graph Disassembly method is that it tries to create disassembly trees rooted
# at every byte-offset of the binary executable, and then invalidate those
# starting byte-offsets for trees in which leaf nodes correspond to invalid
# assembly encodings. Pruning those paths, what is left is but a tree filled
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

# Disassembly.
from capstone import *
from capstone.x86 import *

# Binary file extension.
ext = ".bin"

@click.command()
@click.option('--bin_server', default='127.0.0.1:1234', help='Binary parser server.', show_default=True)
@click.argument('bin_id', type=str)
def client_cmd(bin_server, bin_id):
	# Validate binary ID
	validate_id(bin_id)
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
				sect_data = bin_data[sect.offset:sect.offset+sect.length]
				shingled_disasm(sect, sect_data)

# shingled_disasm determines superset of valid instructions in the given section
# data, returning the valid byte offsets.
def shingled_disasm(sect, sect_data):
	print(sect)
	hexdump.hexdump(sect_data)

	# TODO: make machine architecture configurable.
	cs = Cs(CS_ARCH_X86, CS_MODE_32)
	cs.detail = True

	n = len(sect_data)
	valid = [False] * n
	for i in range(n):
		print("decode: %d / %d" % (i, n))
		addr = sect.addr + i
		valid[i] = is_valid_enc(cs, sect_data[i:], addr)
	# TODO: add info from execution traces.
	visited = [False] * n
	# Queue for decoding.
	q_buf = []
	q_i = 0
	for i in range(n):
		print("validate: %d / %d" % (i, n))
		if visited[i] or not valid[i]:
			continue
		# q_buf[q_i:len(q_buf)] holds the queue of unprocessed instructions
		# reachable from offset i.
		#
		# q_buf[0:q_i] holds the already processed instructions of the queue,
		# reachable from offset i.
		#
		# Together they track the instructions reachable from offset i.
		#
		#q.reset()
		q_buf = []
		q_i = 0
		# Add offset i to queue, and mark as visited.
		visited[i] = True
		#q.push()
		q_buf.append(i)
		prune = False
		#while not q.is_empty():
		while not q_i == len(q_buf):
			#off = q.pop()
			off = q_buf[q_i]
			q_i += 1
			if not valid[off]:
				prune = True
				continue
			addr = sect.addr + off
			inst = decode(cs, sect_data[off:], addr)
			print("0x%x:\t%s\t%s" % (inst.address, inst.mnemonic, inst.op_str))
			if has_fallthrough(inst):
				next_off = off + inst.size
				if next_off < 0 or next_off >= len(sect_data):
					# Skip fallthrough target outside of section data (possible interpretation
					# of data as code at end of segment may have resulted in this operand).
					continue
				if not visited[next_off]:
					visited[next_off] = True
					#q.push(next_off)
					q_buf.append(next_off)
			for branch_addr in branches(inst):
				branch_off = branch_addr - sect.addr
				if branch_off < 0 or branch_off >= len(sect_data):
					# Skip branch target outside of section data (possible interpretation of
					# data as code may have resulted in this operand).
					continue
				if not visited[branch_off]:
					visited[branch_off] = True
					#q.push(branch_off)
					q_buf.append(branch_off)
		if prune:
			# Invalidate all offsets reachable from offset i any of the reachable offsets
			# has an invalid instruction.
			for off in q_buf:
				valid[off] = False
	# Return list of valid instruction offsets.
	return valid

# branches returns the byte offsets of the branches of the given instruction.
def branches(inst):
	opcode = inst.mnemonic
	# Loop terminators.
	if opcode in ['loop', 'loope', 'loopne']:
		operand = inst.operands[0]
	# Conditional jump terminators.
	elif opcode in ['ja', 'jae', 'jb', 'jbe', 'jcxz', 'je', 'jecxz', 'jg', 'jge', 'jl', 'jle', 'jne', 'jno', 'jnp', 'jns', 'jo', 'jp', 'jrcxz', 'js']:
		operand = inst.operands[0]
	# Unconditional jump terminators.
	elif opcode in ['jmp']:
		operand = inst.operands[0]
	# Return terminators.
	elif opcode in ['ret']:
		return []
	# Call instruction (not terminator, but includes branch other than
	# fallthrough).
	elif opcode in ['call']:
		operand = inst.operands[0]
	else:
		return []
	next_addr = inst.address + inst.size
	if operand.type == X86_OP_IMM:
		#return [next_addr + operand.imm]
		return [operand.imm]
	else:
		# TODO: get branch target by symbolic exectuion.
		return []

# decode decodes an instruction at the start of the given data.
def decode(cs, data, addr):
	# TODO: determine max instruction length from cs.
	MAX_INST_LEN = 15
	data = data[:MAX_INST_LEN]
	for inst in cs.disasm(data, addr):
		return inst
	return None

# has_fallthrough reports whether the given instruction has a fall-through
# control flow characteristic.
def has_fallthrough(inst):
	opcode = inst.mnemonic
	return not opcode in ['ret', 'jmp']

# is_valid_enc reports whether the start of the given data encodes a valid
# instruction.
def is_valid_enc(cs, data, addr):
	# TODO: determine max instruction length from cs.
	MAX_INST_LEN = 15
	data = data[:MAX_INST_LEN]
	insts = [inst for inst in cs.disasm(data, addr)]
	return len(insts) > 0

# validate_id validates the given binary ID, terminating the applicatoin if
# invalid.
def validate_id(bin_id):
	if not valid_id(bin_id):
		print('invalid ID; expected lowercase sha256 hash, got %s' % (bin_id))
		sys.exit(1)

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
