from pwn import *
import binascii

# to run to encode and run other programs we can run our go binary
# os.system("./analyser build tmp/scratch.txt tmp/scratch.bin")
# os.system("./analyser analyse code.bin code.dump")

guess = b"n3w_ag3_v1rtu4liz4t1on_" +  bytearray([0x6c, 0x30, 0x6c, 0x3f])
p = process(["./eVMoji", "code.bin"])
out = p.recvuntil("\xf0\x9f\x8f\xb3\xef\xb8\x8f\n")
print(out.decode("utf-8"))
p.send(guess)
out = p.recvall()
print(out.decode("utf-8"))
