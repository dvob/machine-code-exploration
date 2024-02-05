.global _start

.text

sleep:
    enter $32, $0       # reserve 32 byte for two timespec (man timespec) structs
    movq $0, -8(%rbp)
    movq %rax, -16(%rbp) # seconds of first timespec
    movq $0, -24(%rbp)
    movq $0, -32(%rbp)

    mov $35, %rax        # number of nanosleep (man 2 nanosleep) system call
    lea -16(%rbp), %rdi  # address of first timespec struct
    lea -32(%rbp), %rsi  # address of second timespec struct
    syscall

    leave
    ret

_start:
    // wait some time ot observe difference
    // observe page map and page faults with the following commands:
    //   watch -n 1 'pmap $( pgrep a.out )'
    //   watch -n 1 'ps -eo maj_flt,min_flt,cmd -q $( pgrep a.out )'
    mov $5, %rax
    call sleep
    movq $0, -131000(%rsp) # causes the stack to grow one page
    // movq $0, -8400000(%rsp) # causes seg fault because address is higher then ulimts max stack size

    mov $9999, %rax # sleep for a long time to observe memory state
    call sleep

    mov $0, %rdi   # exit code 0
    mov $60, %rax  # system call 60 is exit
    syscall
