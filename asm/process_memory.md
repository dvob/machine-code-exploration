# Memory

CONFIG_PAGE_OFFSET boundary between kernel and userspace
* x86 32bit = 0xC0000000 (upper 1gb, 3gb left for each process)
* x86 64bit = 0xFFFF880000000000

linux device drivers chapter 15 explains the 3 different kind of mappings
* Kernel
  * Kernel Logical Addresses
  * Kernel Virtual Addresses
* User space
  * User Virtual Addresses

A memory page is typically 4k (4096 byte)

$ getconf PAGESIZE

If we try to access an address which is not in the TLB a page fault happens.
That means that the control gets passed to the Kernel and the Kernel has to tell if and where the address is mapped in the physical memory.
a page fault can have three results:
* minor fault: kernel enters the address into the TLB and returns
* major fault: kernel has to load data from disk into memory (mmap, swapping)
* segmentation fault: process has no access to that part of the memory

If the kernel swaps out a page it copies the content of the page to the disk and removes the TLB entry for that page.

Memory is not physically allocated until it is used (lazy allocation).
With mlock syscalls we can force phyisical allocation of memory to avoid later page faults.

To allocate memory from the user space we have to use one of the following syscalls:
* brk
* sbrk
* mmap (with MAP_ANONYMOUS, or MAP_SHARED to share with other processes)

Contrary to what the man page of brk says it does not extend the data section but the heap which is a seperate section.
The \*alloc functions use brk for small allocations and mmap for bigger allocations (see M_MMAP_THRESHOLD in man mallopt).

If we access data near the current stack pointer a new pages is automatically allocated:
https://unix.stackexchange.com/questions/145557/how-does-stack-allocation-work-in-linux

# TODO
run a sleep in a minimal program to observe the pmap of it

# Links
https://www.youtube.com/watch?v=7aONIVSXiJ8

vdso vsyscall: https://0xax.gitbooks.io/linux-insides/content/SysCall/linux-syscall-3.html
virtual system calls: https://lwn.net/Articles/615809/
