#include <stdio.h>
#include <string.h>
#include <sys/mman.h>

int main() {
  unsigned char code[] = {
      0x48, 0xc7, 0xc0, 0x3c, 0x00, 0x00, 0x00, // mov $60, %rax
      0x48, 0xc7, 0xc7, 0x21, 0x00, 0x00, 0x00, // mov $33, %rdi
      0x0f, 0x05                                // syscall
  };

  void (*fn)();
  fn = mmap(NULL, sizeof(code), PROT_EXEC | PROT_READ | PROT_WRITE,
            MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
  memcpy(fn, code, sizeof(code));
  fn();
}
