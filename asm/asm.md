# Assembly Howto
## Registers
| bits | 
| ---- | --- | --- | --- | --- | --- | --- | --- | --- | 
| 8    | AL  | CL  | DL  | BL  | AH  | CH  | DH  | BH  |
| 16   | AX  | CX  | DX  | BX  | SP  | BP  | SI  | DI  |
| 32   | EAX | ECX | EDX | EBX | ESP | EBP | ESI | EDI |
| 64   | RAX | RCX | RDX | RBX | RSP | RBP | RSI | RDI |

* CS, DS, SS, ES: segments (code, data, stack, extra)

* General purpose:
  * RAX
  * RBX
  * RCX
  * RDX
  * RDI
  * RSI
  * RBP (call frame pointer, stack frame)
  * RSP (stack pointer)

* R8 - R15: extra registers (only on 64bit)
* RIP: instruction pointer

## Instructions
* MOV 16bit
* MOVL 32bit (L suffix stands for long)
* MOVQ 64bit (Q suffix stands for quad)

## Function Calls
On 32Bit systems arguments are passed on the stack.
On x86_64 systems the first six arguments are passed with the registers:
* rdi
* rsi
* rdx
* rcx
* r8
* r9

If you have more than 6 arguments the rest are passed via the stack.

The function return value is placed in RAX (and RDX if more bits needed).

The registers %rbx, %rbp, and %r12-15 are considered callee-saved registers.
This means they are not freely available to use in a function;
If a function wants to use them, it must first save them (by pushing them on the stack), and then restore them at the end.
All other registers are considered freely available (the caller must save them before the call if it needs their values after the call).

On function calls the stack has to be aligned properly. In case of System V x86_64 the alignment is 16byte.

You can observe this if you do something like this:
```c
int do_stuff() {
    int v;
    v = foo();
    return 0;
}
```

The assembly code of such a function would start like this:
```asm
do_stuff:
    push %rbp
    mov %rsp, %rbp
    sub $0x10, %rsp
```
With sub 16byte (0x10=16) are subtracted from the stack pointer even though the local variable v has only a size of 4bytes.

See the [`func.s`](./func.s) example for more informatino about function calls.

## Syscalls
In x86_32 a programm runs the interrupt instruction with 128 (0x80) as argument. In x86_64 syscall is an own instruction. See https://serverfault.com/a/880814.

Also in 32bit and 64bit systems syscalls do not use the same number:
* 64Bit: https://github.com/torvalds/linux/blob/master/arch/x86/entry/syscalls/syscall_64.tbl
* 32Bit: https://github.com/torvalds/linux/blob/master/arch/x86/entry/syscalls/syscall_32.tbl

See also the following link for a nice overview: https://filippo.io/linux-syscall-table/


To call a syscall we write the syscall number to %rax and the arguments to the appropriate registers (see function call).
The result will be written to %rax.

# Debug
Run the debugger:
```
gcc -g -c hello.s && ld hello.o
gdb ./a.out
```

In the debugger:
```
# break point on _start
b _start

# run with stdin
run <input.txt

# or just
run

# then you can step thorugh the instructions with
s

# print a string in memory
x/s $rbp - 80

# inspect registers
i r
```

## Read arguments
```
b _start
run arg1 arg2 arg3

# argc is on $rsp
x/1u $rsp

# argv is on $rsp+8
# get address of first argument
x/1a $rsp+8

# then copy address of output and print the string
x/1s 0x7fffffffea2c
```

# Links
* Upper bits of a 64bit register (e.g. rax) get zeroed if you use eax: https://stackoverflow.com/questions/11177137/why-do-x86-64-instructions-on-32-bit-registers-zero-the-upper-part-of-the-full-6

* GNU Assembler Doc: https://sourceware.org/binutils/docs/as/index.html

* What happens before main(): https://embeddedartistry.com/blog/2019/04/08/a-general-overview-of-what-happens-before-main/

* GAS examples: 
  * https://codedocs.org/what-is/gnu-assembler
  * https://cs.lmu.edu/~ray/notes/gasexamples/

* Function calls: https://zhu45.org/posts/2017/Jul/30/understanding-how-function-call-works/

* Assembly cheat sheet: https://cs.brown.edu/courses/cs033/docs/guides/x64_cheatsheet.pdf

* Registers: https://en.wikibooks.org/wiki/X86_Assembly/16,_32,_and_64_Bits#Registers
* Difference 16bit, 32bit 64bit: https://www.cs.nmsu.edu/~jcook/posts/basic-x86-64-assembly/  

* Hello World: https://gist.github.com/carloscarcamo/6833d19b726af698e62b
* Assembly Tutorial: https://riptutorial.com/assembly
* Linux System Call Table for x86 64: http://blog.rchapman.org/posts/Linux_System_Call_Table_for_x86_64/
* x86 Manual: http://www.intel.com/content/dam/www/public/us/en/documents/manuals/64-ia-32-architectures-software-developer-instruction-set-reference-manual-325383.pdf
* x86 Assembly Guide: https://www.cs.virginia.edu/~evans/cs216/guides/x86.html

* System V Application Binary Interface: https://refspecs.linuxbase.org/elf/x86_64-abi-0.99.pdf

* Explains various features in NASM and GAS: https://developer.ibm.com/articles/l-gas-nasm/

* RISC-V https://risc-v.guru/

* Run assembly from C: https://gcc.gnu.org/onlinedocs/gcc/Simple-Constraints.html#Simple-Constraints

# Books
* Assembly Language Step-by-Step, Jeff Duntemann: http://duntemann.com/assembly.html
