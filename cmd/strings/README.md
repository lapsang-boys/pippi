
# strings

```bash
$ (strings/) cargo build
$ (strings/) ./target/debug/strings

# From project root
$ (/) grpc_cli call localhost:1234 ExtractStrings "elf_path: '/usr/bin/ls'" --proto_path proto/ --protofiles strings.proto
```
