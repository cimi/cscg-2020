from pwn import *
import binascii

prefix = b"n3w_ag3_v1rtu4liz4t1on_"


# in program: 01101100001100000110110000111111

# 0011 0110 0000 1100 0011 0110 1111 1100
# 3    6    0    c    3    6    f    c
# 🎉🎉 Solution: 0x360c36fc

# 0x360c36cf => 0110 1100 0011 0000 0110 1100 0011 1111 | 6c 30 6c 3f | f3 c6 03 c6
# 0xf3c603c6 => 1100 0110 0000 0011 1100 0110 1111 0011 | c6 03 c6 f3
guesses = [
  bytearray([0xfc, 0x36, 0x0c, 0x36]),
  bytearray([0x3f, 0x6c, 0x30, 0x6c]),
 ,
]

# to run to encode and run other programs we can run our go binary
# os.system("./emoji-analyzer build tmp/scratch.txt tmp/scratch.bin")

guess = prefix +  bytearray([0x6c, 0x30, 0x6c, 0x3f])
p = process(["./eVMoji", "code.bin"])
p.recvuntil("\xf0\x9f\x8f\xb3\xef\xb8\x8f\n") # white flag emoji
p.send(guess)
out = p.recvall() # we have the leaked value in the output
print(out)
