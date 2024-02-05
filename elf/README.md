# Executable and Linking Format

The specification of the ELF format was first published in the System V ABI specification and later in the Tool Interface Standard [^1].

Specification:
* System V ABI specifications
  * [System V ABI - Base](http://www.sco.com/developers/devspecs/gabi41.pdf)
  * [System V ABI - AMD64 Supplement](https://refspecs.linuxbase.org/elf/x86_64-abi-0.99.pdf)
  * also see https://wiki.osdev.org/System_V_ABI#Documents for more System V specifications
* [Tool Interface Standard (TIS) Executable and Linking Format (ELF) Specification Version 1.2](https://refspecs.linuxbase.org/elf/elf.pdf)

# Format

* Header
* Program headers (segments): where and how to put the code from the file into memory on process start.
* Section headers
  * various information (can point into segments)
  * required for the linker
  * debug symbols
  * not required for runtime

Common sections are:
* `.text`: executable code
* global and static variables
  * `.data`: data with initial value (`int i = 5;`);
  * `.bss`: only reserved space (`int b;`);
* `.rodata`: constant (read only) data. e.g. string or byte arrays.
* `.symtab`: symbol table for functions, global variables, constants, etc.
* `.strtab` and `.shstrtab`: strings describing the names of other ELF structures e.g. in `.symtab` or even section names

## Links
Read and run code from an ELF file using Rust (18 parts)
* https://fasterthanli.me/series/making-our-own-executable-packer/part-1

Introduction into ELF
* https://blog.k3170makan.com/2018/09/introduction-to-elf-format-elf-header.html

ELF, PIC, GOT:
* https://www.linuxjournal.com/article/1059

How to load and run code form a .o file:
* https://blog.cloudflare.com/how-to-execute-an-object-file-part-1/
* https://blog.cloudflare.com/how-to-execute-an-object-file-part-2/
* https://blog.cloudflare.com/how-to-execute-an-object-file-part-3/

Minimal ELF files:


[^1]: https://en.wikipedia.org/wiki/Executable_and_Linkable_Format
