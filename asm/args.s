.global _start

.text

print_newline:
    enter $0, $0
    push $10
    mov $1, %rax            # system call 1 is write
    mov $1, %rdi            # file handler 1 is stdout
    lea -8(%rbp), %rsi      # address of string to output
    mov $1, %rdx
    syscall
    leave
    ret


// get the length of a zero terminated string
get_len:
    enter $8, $0        # same as above

    movq $0, -8(%rbp)   # initialize len counter to zero

loop:
    // jump to end of function if NULL char is found
    cmpb $0, (%rdi)
    je end

    inc %rdi        # increase char pointer
    incq -8(%rbp)   # increase pointer
    jmp loop

end:
    mov -8(%rbp), %rax
    leave
    ret


_start:
    mov %rsp, %rbp
    movq (%rsp), %r11      # argc
    movq %r11, -8(%rbp)    # argc

    lea 8(%rsp), %r11     # argv
    mov %r11, -16(%rbp)   # argv
    movq $0, -24(%rbp)       # counter = 0
    sub $32, %rsp

args:
    // counter == argc
    mov -8(%rbp), %r11
    cmpq -24(%rbp), %r11
    je exit

    // get len
    mov -16(%rbp), %rdi
    mov (%rdi), %rdi
    call get_len

    mov %rax, %rdx          # you len from get_len as len argument for write
    mov $1, %rax            # system call 1 is write
    mov $1, %rdi            # file handler 1 is stdout

    mov -16(%rbp), %rsi      # address of string to output
    mov (%rsi), %rsi
    syscall

    call print_newline

    cmp $0, %rax
    jl errexit

    incq -24(%rbp)      # increase counter
    addq $8, -16(%rbp)   # jump to next argument
    jmp args

    // mov $my_str, %rdi
    // call get_len
errexit:
    mov $66, %rdi
exit:
    mov $0, %rdi
    mov $60, %rax           # system call 60 is exit
    syscall
