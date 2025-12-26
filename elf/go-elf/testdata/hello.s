.section .text
.globl _start

hello_msg:
    .ascii "Hello\n\0"

_start:
    # write(1, hello_msg, 5)
    mov $1, %rax        # syscall number for write
    mov $1, %rdi        # fd = stdout
    lea hello_msg(%rip), %rsi  # pointer to "Hello"
    mov $6, %rdx        # length (without null terminator)
    syscall

    # exit(33)
    mov $60, %rax       # syscall number for exit
    mov $33, %rdi       # exit code
    syscall
