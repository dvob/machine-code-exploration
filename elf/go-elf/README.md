# go-elf

Read and create an ELF files.

```
go build ./cmd/elf-debug
```

Print information about an ELF file:
```
( cd testdata; gcc main.c )

./elf-debug read testdata/a.out
```

Write an ELF file:
```
./elf-debug write

./output.elf
echo $?
```
