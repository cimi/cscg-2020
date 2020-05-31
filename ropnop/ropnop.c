// clang -fno-stack-protector ./ropnop.c -o ropnop

#include <stdio.h>
#include <unistd.h>
#include <sys/mman.h>

extern unsigned char __executable_start;
extern unsigned char etext;


void init_buffering() {
	setvbuf(stdout, NULL, _IONBF, 0);
	setvbuf(stdin, NULL, _IONBF, 0);
	setvbuf(stderr, NULL, _IONBF, 0);
}


void gadget_shop() {
	// look at all these cool gadgets
	__asm__("syscall; ret");
	__asm__("pop %rax; ret");
	__asm__("pop %rdi; ret");
	__asm__("pop %rsi; ret");
	__asm__("pop %rdx; ret");
}

void ropnop() {
	unsigned char *start = &__executable_start;
	unsigned char *end = &etext;
	printf("[defusing returns] start: %p - end: %p\n", start, end);
	mprotect(start, end-start, PROT_READ|PROT_WRITE|PROT_EXEC);
	unsigned char *p = start;
	while (p != end) {
		// if we encounter a ret instruction, replace it with nop!
		if (*p == 0xc3)
			*p = 0x90;
		p++;
	}
}

int main(void) {
	init_buffering();
	ropnop();
	int* buffer = (int*)&buffer;
	read(0, buffer, 0x1337);
	return 0;
}
