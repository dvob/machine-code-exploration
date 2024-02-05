# -------------------------------------------------------------
#
.global _start

.text

_start:
    mov $1, %rax            # system call 1 is write
    mov $1, %rdi            # file handler 1 is stdout
    mov $message, %rsi      # address of string to output
    mov $13, %rdx           # number of bytes
    syscall

    # exit(0)
    mov $60, %rax           # system call 60 is exit
    mov $0, %rdi
    syscall

message:
    .ascii "Hello, world\n"
