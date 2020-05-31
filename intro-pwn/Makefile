CC=gcc
CFLAGS=

all: pwn1 pwn2 pwn3

pwn1: pwn1/pwn1
	$(CC) -fno-stack-protector pwn1/pwn1.c -o pwn1/pwn1
	
pwn2: pwn2/pwn2
	$(CC) -D PASSWORD=\"`cat pwn1/flag`\" pwn2/pwn2.c -o pwn2/pwn2 

pwn3: pwn3/pwn3
	$(CC) -D PASSWORD=\"`cat pwn2/flag`\" pwn3/pwn3.c -o pwn3/pwn3

clean:
	rm pwn1/pwn1 pwn2/pwn2 pwn3/pwn3