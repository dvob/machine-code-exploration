// also see https://en.wikipedia.org/wiki/Function_prologue_and_epilogue

.global _start

add_one:
    // prologue (setup stack frame)
    push %rbp           # safe base pointer of caller
    mov %rsp, %rbp      # set our own base pointer based on the stack pointer
    sub $8, %rsp        # allocate/reserve space on stack

    // perform calculations on stack
    mov %rdi, -8(%rbp)  # write first argument to stack
    addq $1, -8(%rbp)   # add one to the value on the stack

    // prepare return value
    mov -8(%rbp), %rax

    // epilogue
    mov %rbp, %rsp      # restore stack pointer => inverse of alloc
    pop %rbp            # restore base pointer of caller
    ret                 # pops old instruction pointer from stack and jumps to it => inverse of call

// short form
add_one_short:
    // prologue
    enter $8, $0        # same as above

    // perform calculations on stack
    mov %rdi, -8(%rbp)  # write first argument to stack
    addq $1, -8(%rbp)   # add one to the value on the stack

    // prepare return value
    mov -8(%rbp), %rax

    // epilogue
    leave
    ret

_start:
    // call add_one
    mov $1, %rdi  # prepare first argument
    call add_one  # pushes current instruction pointer (rip) + 1 to stack and jumps to address of add_one

    // call add_one_short
    mov %rax, %rdi  # use result from previous call as first argument
    call add_one_short

    // exit
    mov %rax, %rdi # use result of add_two as parameter to exit()
    mov $60, %rax
    syscall
