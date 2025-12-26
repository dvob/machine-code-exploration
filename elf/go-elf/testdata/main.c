#include <stdio.h>

// initialized to hex value 0xbeef
int counter = 48879;

int main() {
  counter++;
  printf("Hello world! %d", counter);
}
