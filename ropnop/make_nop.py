src = open("ropnop", "rb")
dst = open("nopnop", "wb")

try:
  byte = src.read(1)
  while byte != "":
    # Do stuff with byte.
    if byte == chr(0xc3):
      dst.write('\x90')
    else:
      dst.write(byte)

    byte = src.read(1)        
finally:
  src.close()
  dst.close()



# 0x7f842f1 3ed70: "printf"
# 0x7f842f1 61d50: "setvbuf"
# 0x7f842f1 eb260: "read"
# 0x7f842f1 f5bb0: "mprotect"

# printf@got.plt>
# read@got.plt
# setvbuf@got.plt
# mprotect@got.plt



# pwndbg> rop --grep ret -- --range 0x55b8fb8d4000-0x55b8fb8d6000
# Saved corefile /tmp/tmpu4pqrf7y
# 0x000055b8fb8d535c : add byte ptr [rax], al ; add byte ptr [rax], al ; endbr64 ; ret
# 0x000055b8fb8d535e : add byte ptr [rax], al ; endbr64 ; ret
# 0x000055b8fb8d5297 : add esp, 0x20 ; pop rbp ; ret
# 0x000055b8fb8d5371 : add esp, 8 ; ret
# 0x000055b8fb8d5296 : add rsp, 0x20 ; pop rbp ; ret
# 0x000055b8fb8d5370 : add rsp, 8 ; ret
# 0x000055b8fb8d5363 : cli ; ret
# 0x000055b8fb8d536b : cli ; sub rsp, 8 ; add rsp, 8 ; ret
# 0x000055b8fb8d5360 : endbr64 ; ret
# 0x000055b8fb8d533c : fisttp word ptr [rax - 0x7d] ; ret
# 0x000055b8fb8d52da : mov eax, ecx ; add rsp, 0x20 ; pop rbp ; ret
# 0x000055b8fb8d534c : pop r12 ; pop r13 ; pop r14 ; pop r15 ; ret
# 0x000055b8fb8d534e : pop r13 ; pop r14 ; pop r15 ; ret
# 0x000055b8fb8d5350 : pop r14 ; pop r15 ; ret
# 0x000055b8fb8d5352 : pop r15 ; ret
# 0x000055b8fb8d534b : pop rbp ; pop r12 ; pop r13 ; pop r14 ; pop r15 ; ret
# 0x000055b8fb8d534f : pop rbp ; pop r14 ; pop r15 ; ret
# 0x000055b8fb8d529a : pop rbp ; ret
# 0x000055b8fb8d5353 : pop rdi ; ret
# 0x000055b8fb8d5351 : pop rsi ; pop r15 ; ret
# 0x000055b8fb8d534d : pop rsp ; pop r13 ; pop r14 ; pop r15 ; ret
# 0x000055b8fb8d529b : ret
# 0x000055b8fb8d5121 : retf 0x2e
# 0x000055b8fb8d5062 : retf 0x2f
# 0x000055b8fb8d536d : sub esp, 8 ; add rsp, 8 ; ret
# 0x000055b8fb8d536c : sub rsp, 8 ; add rsp, 8 ; ret



# rop -- --range 0x563384bf2000-0x563384bf4000
# Saved corefile /tmp/tmp1zi7nb2f
# Gadgets information
# ============================================================
# 0x0000563384bf3057 : add al, byte ptr [rax] ; add byte ptr [rax], al ; jmp 0x563384bf3024
# 0x0000563384bf31de : add byte ptr [rax - 0x1a76b7ab], dl ; syscall
# 0x0000563384bf31dc : add byte ptr [rax], al ; add byte ptr [rax - 0x1a76b7ab], dl ; syscall
# 0x0000563384bf335c : add byte ptr [rax], al ; add byte ptr [rax], al ; endbr64 ; ret
# 0x0000563384bf3037 : add byte ptr [rax], al ; add byte ptr [rax], al ; jmp 0x563384bf3024
# 0x0000563384bf31db : add byte ptr [rax], al ; add byte ptr [rax], al ; nop ; push rbp ; mov rbp, rsp ; syscall
# 0x0000563384bf335e : add byte ptr [rax], al ; endbr64 ; ret
# 0x0000563384bf3039 : add byte ptr [rax], al ; jmp 0x563384bf3022
# 0x0000563384bf328b : add byte ptr [rax], al ; mov qword ptr [rbp - 0x18], rax ; jmp 0x563384bf3261
# 0x0000563384bf31dd : add byte ptr [rax], al ; nop ; push rbp ; mov rbp, rsp ; syscall
# 0x0000563384bf3034 : add byte ptr [rax], al ; push 0 ; jmp 0x563384bf3027
# 0x0000563384bf3044 : add byte ptr [rax], al ; push 1 ; jmp 0x563384bf3027
# 0x0000563384bf3054 : add byte ptr [rax], al ; push 2 ; jmp 0x563384bf3027
# 0x0000563384bf3064 : add byte ptr [rax], al ; push 3 ; jmp 0x563384bf3027
# 0x0000563384bf300d : add byte ptr [rax], al ; test rax, rax ; je 0x563384bf301d ; call rax
# 0x0000563384bf30b8 : add byte ptr [rax], al ; test rax, rax ; je 0x563384bf30cf ; jmp rax
# 0x0000563384bf30f9 : add byte ptr [rax], al ; test rax, rax ; je 0x563384bf310f ; jmp rax
# 0x0000563384bf30f8 : add byte ptr cs:[rax], al ; test rax, rax ; je 0x563384bf3110 ; jmp rax
# 0x0000563384bf3047 : add dword ptr [rax], eax ; add byte ptr [rax], al ; jmp 0x563384bf3024
# 0x0000563384bf3289 : add dword ptr [rax], eax ; add byte ptr [rax], al ; mov qword ptr [rbp - 0x18], rax ; jmp 0x563384bf3263
# 0x0000563384bf3288 : add eax, 1 ; mov qword ptr [rbp - 0x18], rax ; jmp 0x563384bf3264
# 0x0000563384bf3067 : add eax, dword ptr [rax] ; add byte ptr [rax], al ; jmp 0x563384bf3024
# 0x0000563384bf3297 : add esp, 0x20 ; pop rbp ; ret
# 0x0000563384bf3371 : add esp, 8 ; ret
# 0x0000563384bf3296 : add rsp, 0x20 ; pop rbp ; ret
# 0x0000563384bf3370 : add rsp, 8 ; ret
# 0x0000563384bf3014 : call rax
# 0x0000563384bf3163 : cli ; jmp 0x563384bf30d1
# 0x0000563384bf3363 : cli ; ret
# 0x0000563384bf336b : cli ; sub rsp, 8 ; add rsp, 8 ; ret
# 0x0000563384bf3160 : endbr64 ; jmp 0x563384bf30d4
# 0x0000563384bf3360 : endbr64 ; ret
# 0x0000563384bf333c : fisttp word ptr [rax - 0x7d] ; ret
# 0x0000563384bf3042 : fisubr dword ptr [rdi] ; add byte ptr [rax], al ; push 1 ; jmp 0x563384bf3029
# 0x0000563384bf30f7 : in eax, dx ; add byte ptr cs:[rax], al ; test rax, rax ; je 0x563384bf3111 ; jmp rax
# 0x0000563384bf3012 : je 0x563384bf3018 ; call rax
# 0x0000563384bf30bd : je 0x563384bf30ca ; jmp rax
# 0x0000563384bf30fe : je 0x563384bf310a ; jmp rax
# 0x0000563384bf303b : jmp 0x563384bf3020
# 0x0000563384bf3164 : jmp 0x563384bf30d0
# 0x0000563384bf3291 : jmp 0x563384bf325b
# 0x0000563384bf30bf : jmp rax
# 0x0000563384bf3032 : loop 0x563384bf306c ; add byte ptr [rax], al ; push 0 ; jmp 0x563384bf3029
# 0x0000563384bf328e : mov dword ptr [rbp - 0x18], eax ; jmp 0x563384bf325e
# 0x0000563384bf32da : mov eax, ecx ; add rsp, 0x20 ; pop rbp ; ret
# 0x0000563384bf31e2 : mov ebp, esp ; syscall
# 0x0000563384bf328d : mov qword ptr [rbp - 0x18], rax ; jmp 0x563384bf325f
# 0x0000563384bf31e1 : mov rbp, rsp ; syscall
# 0x0000563384bf31df : nop ; push rbp ; mov rbp, rsp ; syscall
# 0x0000563384bf315c : nop dword ptr [rax] ; endbr64 ; jmp 0x563384bf30d8
# 0x0000563384bf334c : pop r12 ; pop r13 ; pop r14 ; pop r15 ; ret
# 0x0000563384bf334e : pop r13 ; pop r14 ; pop r15 ; ret
# 0x0000563384bf3350 : pop r14 ; pop r15 ; ret
# 0x0000563384bf3352 : pop r15 ; ret
# 0x0000563384bf334b : pop rbp ; pop r12 ; pop r13 ; pop r14 ; pop r15 ; ret
# 0x0000563384bf334f : pop rbp ; pop r14 ; pop r15 ; ret
# 0x0000563384bf329a : pop rbp ; ret
# 0x0000563384bf3353 : pop rdi ; ret
# 0x0000563384bf3351 : pop rsi ; pop r15 ; ret
# 0x0000563384bf334d : pop rsp ; pop r13 ; pop r14 ; pop r15 ; ret
# 0x0000563384bf3036 : push 0 ; jmp 0x563384bf3025
# 0x0000563384bf3046 : push 1 ; jmp 0x563384bf3025
# 0x0000563384bf3056 : push 2 ; jmp 0x563384bf3025
# 0x0000563384bf3066 : push 3 ; jmp 0x563384bf3025
# 0x0000563384bf31e0 : push rbp ; mov rbp, rsp ; syscall
# 0x0000563384bf329b : ret
# 0x0000563384bf3121 : retf 0x2e
# 0x0000563384bf3062 : retf 0x2f
# 0x0000563384bf3052 : shr byte ptr [rdi], cl ; add byte ptr [rax], al ; push 2 ; jmp 0x563384bf3029
# 0x0000563384bf300b : shr dword ptr [rdi], 1 ; add byte ptr [rax], al ; test rax, rax ; je 0x563384bf301f ; call rax
# 0x0000563384bf336d : sub esp, 8 ; add rsp, 8 ; ret
# 0x0000563384bf336c : sub rsp, 8 ; add rsp, 8 ; ret
# 0x0000563384bf31e4 : syscall
# 0x0000563384bf3010 : test eax, eax ; je 0x563384bf301a ; call rax
# 0x0000563384bf30bb : test eax, eax ; je 0x563384bf30cc ; jmp rax
# 0x0000563384bf30fc : test eax, eax ; je 0x563384bf310c ; jmp rax
# 0x0000563384bf300f : test rax, rax ; je 0x563384bf301b ; call rax
# 0x0000563384bf30ba : test rax, rax ; je 0x563384bf30cd ; jmp rax
# 0x0000563384bf30fb : test rax, rax ; je 0x563384bf310d ; jmp rax

# Unique gadgets found: 79
