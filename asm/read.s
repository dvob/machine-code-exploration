.global _start

.data
input: .asciz "input.txt"

.text
_start:
    // function prolog
    push %rbp
    mov %rsp, %rbp
    sub $528, %rsp  # 512 Byte for data + 16 byte for variables

    // write
    mov $1, %rax            # system call 1 is write
    mov $1, %rdi            # file handler 1 is stdout
    mov $input, %rsi      # address of string to output
    mov $9, %rdx           # number of bytes
    syscall

    // open input.txt
    mov $2, %rax
    mov $input, %rdi
    mov $0, %rsi # O_RDONLY
    syscall
    mov %rax, -8(%rbp)

    // check for open errors
    mov -8(%rbp), %rdi # set bad exitcode
    cmpq $0, -8(%rbp) 
    jl exit

loop:
    // read from open filedescriptor
    mov $0, %rax            # read syscall
    mov -8(%rbp), %rdi      # file descriptor
    lea -528(%rbp), %rsi     # write to our stack
    mov $512, %rdx           # buffer size
    syscall
    mov %rax, -16(%rbp)

    // check for read errors
    mov -16(%rbp), %rdi # set bad exitcode
    cmpq $0, -16(%rbp) 
    jl exit

    cmpq $0, -16(%rbp) 
    je loop_end

    // write
    mov $1, %rax         # system call 1 is write
    mov $1, %rdi   # file descriptor
    lea -528(%rbp), %rsi # address of buffer
    mov -16(%rbp), %rdx  # number read in previous call to read
    syscall
    mov %rax, -16(%rbp)

    // check for write error
    mov -16(%rbp), %rdi # set bad exitcode
    cmpq $0, -16(%rbp) 
    jl exit

    jmp loop

loop_end:
    
    mov $0, %rdi

exit:
    mov $60, %rax           # system call 60 is exit
    #mov $rax, %rdi
    syscall
