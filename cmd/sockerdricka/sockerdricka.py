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
