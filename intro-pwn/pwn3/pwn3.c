#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <signal.h>
#include <string.h>

#ifndef PASSWORD
    #define PASSWORD "CSCG{FLAG_FROM_STAGE_2}"
#endif

// pwn3: gcc pwn3.c -o pwn3

// --------------------------------------------------- SETUP

void ignore_me_init_buffering() {
	setvbuf(stdout, NULL, _IONBF, 0);
	setvbuf(stdin, NULL, _IONBF, 0);
	setvbuf(stderr, NULL, _IONBF, 0);
}

void kill_on_timeout(int sig) {
  if (sig == SIGALRM) {
  	printf("[!] Anti DoS Signal. Patch me out for testing.");
    _exit(0);
  }
}

void ignore_me_init_signal() {
	signal(SIGALRM, kill_on_timeout);
	alarm(60);
}

// just a safe alternative to gets()
size_t read_input(int fd, char *buf, size_t size) {
  size_t i;
  for (i = 0; i < size-1; ++i) {
    char c;
    if (read(fd, &c, 1) <= 0) {
      _exit(0);
    }
    if (c == '\n') {
      break;
    }
    buf[i] = c;
  }
  buf[i] = '\0';
  return i;
}

// --------------------------------------------------- MENU

void WINgardium_leviosa() {
    printf("They has discovered our secret, Nagini.\n");
    printf("It makes us vulnerable.\n");
    printf("We must deploy all our forces now to find them.\n");
    // system("/bin/sh") it's not that easy anymore.
}

void check_password_stage2() {
    char read_buf[0xff];
    printf("Enter the password of stage 2:\n");
    memset(read_buf, 0, sizeof(read_buf));
    read_input(0, read_buf, sizeof(read_buf));
    if(strcmp(read_buf, PASSWORD) != 0) {
        printf("-10 Points for Gryffindor!\n");
        _exit(0);
    } else {
        printf("+10 Points for Gryffindor!");
    }
}

void welcome() {
    char read_buf[0xff];
    printf("Enter your witch name:\n");
    gets(read_buf);
    printf("┌───────────────────────┐\n");
    printf("│ You are a Gryffindor! │\n");
    printf("└───────────────────────┘\n");
    printf(read_buf);
}

void AAAAAAAA() {
    char read_buf[0xff];
    
    printf(" enter your magic spell:\n");
    gets(read_buf);
    if(strcmp(read_buf, "Expelliarmus") == 0) {
        printf("~ Protego!\n");
    } else {
        printf("-10 Points for Gryffindor!\n");
        _exit(0);
    }
}

// --------------------------------------------------- MAIN

void main(int argc, char* argv[]) {
	  ignore_me_init_buffering();
	  ignore_me_init_signal();

    check_password_stage2();

    welcome();
    AAAAAAAA();	
}
