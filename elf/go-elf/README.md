# go-elf

Read and create an ELF file.

Build:
```shell
# build a.out
gcc -nostdlib -static ../../exec/main.s

# build go-elf
go build
```

```shell
# creates output.elf
./go-elf a.out
```

```shellp
./output.elf
echo $?
```
