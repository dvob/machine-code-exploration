.global _start

.bss

.lcomm input 12000 # allocate n bytes

.data
my_str: .asciz "14" # zero terminated string

.text
//
// get the length of a zero terminated string
//
get_len:
    enter $8, $0        # same as above

    movq $0, -8(%rbp)   # initialize len counter to zero

  get_len_loop:
    // jump to end of function if NULL char is found
    cmpb $0, (%rdi)
    je get_len_end

    inc %rdi        # increase char pointer
    incq -8(%rbp)   # increase pointer
    jmp get_len_loop

  get_len_end:
    mov -8(%rbp), %rax
    leave
    ret

//
// read into
//    args:
//      - fd file descriptor
//      - destination *char
//      - max len
//
read_into:
    enter $32, $0        # same as above

    // move arguments to stack
    mov %rdi, -8(%rbp)  # filedescriptor
    mov %rsi, -16(%rbp) # pointer to memory
    mov %rdx, -24(%rbp) # max len
    movq $0, -32(%rbp)  # initialize len counter to zero

    mov $0, %rax            # read syscall
    mov -8(%rbp), %rdi      # file descriptor
    mov -16(%rbp), %rsi     # write to our stack
    mov -24(%rbp), %rdx
    syscall
    // rax already set from syscall
    leave
    ret

// find_char:
//     enter $0, $0        # same as above
// 
//     mov %rdi        # pointer
//     movq $0, %r8    # counter
// 
//   find_char_comp:
//     cmpb $'0', %(rdi)
//     jb find_char_exit
// 
//     cmpb $'9', %(rdi)
//     ja find_char_exit
// 
//     inc %r8
//     inc %rdi
//     jmp find_char_comp
// 
//   find_char_exit:
//     mov %r8, %rax
//     leave
//     ret

num_len:
    enter $0, $0

    movq $0, %rax

  num_len_loop:
    cmpb $'0', (%rdi)
    jb num_len_exit

    cmpb $'9', (%rdi)
    ja num_len_exit

    inc %rax
    inc %rdi

    jmp num_len_loop

  num_len_exit:
    leave
    ret

// read numeric chars (0-9) and return number
atoi:
    enter $16, $0

    // move arguments to stack
    // mov %rdi, -8(%rbp)  # ptr
    // mov $0, -16(%rbp) # number

    mov $0, %rax

    // first iteration
    cmpb $'0', (%rdi)
    jb atoi_exit
 
    cmpb $'9', (%rdi)
    ja atoi_exit

    movq $0, %rsi
    movb (%rdi), %sil
    subq $'0', %rsi
    addq %rsi, %rax

  atoi_loop:
    inc %rdi
    cmpb $'0', (%rdi)
    jb atoi_exit
 
    cmpb $'9', (%rdi)
    ja atoi_exit

    movq $10, %r8
    mul %r8
    movb (%rdi), %sil
    subq $'0', %rsi
    addq %rsi, %rax

    jmp atoi_loop

  atoi_exit:
    leave
    ret

_start:
    // open argv[1]
    mov $2, %rax
    mov 16(%rsp), %rdi
    mov $0, %rsi # O_RDONLY
    syscall

    // check result of open
    mov %rax, %rdi
    cmp $0, %rdi
    jl exit_start

    // read input into global memory input in .bss section
    // fd already moved to rsi above
    mov $input, %rsi
    mov $12000, %rdx
    call read_into

    // check result of read_into
    mov %rax, %rdi
    cmp $0, %rdi
    jl exit_start

    mov $my_str, %rdi
    call num_len
    mov %rax, %rdi
    jmp exit_start


    mov $0, %rdi  # exit code 0
  exit_start:
    mov $60, %rax           # system call 60 is exit
    syscall
