# Execute

Execute machine code as part of a "regular" program:
* mmap some memory space which is writable and executable
* copy the machine code into it
* turn the pointer to the mmapped space into a function
* run it

Such a technique would, for example, be used in a JIT compiler.

Examples:
* [C](./c-exec/)
* [Rust](./rs-exec/)
* [Go](./go-exec/)

## Minimal example
As a minimal example of machine code, we use a program that simply exits with exit code 33 ([`main.s`](./main.s)).
```
gcc -nostdlib -static main.s
```

```
objcopy -O binary --only-section=.text a.out my.bin
```

## Disassemble
* Disassemble machine code that is not part of an ELF file:
```
echo -e -n '\x48\xc7\xc0\x3c\x00\x00\x00\x48\xc7\xc7\x21\x00\x00\x00\x0f\x05' > my.bin
objdump -b binary -D -m i386:x86-64  my.bin
```
