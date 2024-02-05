.global _start

.text

_start:
    mov $42, %rdi  # exit code 0
    mov $60, %rax           # system call 60 is exit
    syscall
