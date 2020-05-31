from pwn import *
import binascii

prefix = b"n3w_ag3_v1rtu4liz4t1on_"
os.system("./emoji build tmp/scratch.txt tmp/scratch.bin")

p128 = bytearray([0x20, 0x83, 0xb8, 0xed])
p136 = bytearray([0x5e, 0x84, 0x0e, 0xf4])
p140 = bytearray([0xff, 0xff, 0xff, 0xff])
p100 = bytearray([0x00, 0x00, 0x00, 0x00])

def process_output(out):
  return out.decode('utf-8').strip().split("Gotta go")[0].replace("3", "1").replace("n", "0").replace("w", "|")


# in program: 01101100001100000110110000111111

# 0011 0110 0000 1100 0011 0110 1111 1100
# 3    6    0    c    3    6    f    c
# 🎉🎉 Solution: 0x360c36fc

# 0x360c36cf => 0110 1100 0011 0000 0110 1100 0011 1111 | 6c 30 6c 3f | f3 c6 03 c6
# 0xf3c603c6 => 1100 0110 0000 0011 1100 0110 1111 0011 | c6 03 c6 f3
guesses = [
  bytearray([0xfc, 0x36, 0x0c, 0x36]),
  bytearray([0x3f, 0x6c, 0x30, 0x6c]),
  bytearray([0x6c, 0x30, 0x6c, 0x3f]),
  # bytearray([0x0f, 0x00, 0x00, 0x00]),
  # bytearray([0x00, 0xf0, 0x00, 0x00]),
  # bytearray([0x00, 0x0f, 0x00, 0x00]),
  # bytearray([0x00, 0x00, 0xf0, 0x00]),
  # bytearray([0x00, 0x00, 0x0f, 0x00]),
  # bytearray([0x00, 0x00, 0x00, 0xf0]),
  # bytearray([0x00, 0x00, 0x00, 0x0f]),
]

# f0000000 00000000000000000000000011110000 - 0->6
# 0f000000 00000000000000000000000000001111 - 1->7
# 00f00000 00000000000000001111000000000000 - 2->4
# 000f0000 00000000000000000000111100000000 - 3->5
# 0000f000 00000000111100000000000000000000 - 4->2
# 00000f00 00000000000011110000000000000000 - 5->3
# 000000f0 11110000000000000000000000000000 - 6->0
# 0000000f 00001111000000000000000000000000 - 7->1

# 0011 0110 0000 1100 0011 0110 1111 1100
# 0011 0110 0000 1100 0011 0110 1111 1100
# 6->3 7->6 4->0 5->c 2->3 3->6 0->f 1->c

# 0011 1111 0110 1100 0011 0000 0110 1100
# 6->3 7->f 4->6 5->c 2->3 3->0 0->6 1->c

# 0011 1111 0110 1100 0011 0000 0110 1100

# pointer = bytearray([0x3f, 0x6c, 0x30, 0x6c])
# for i in range(255, 256):
#   for j in range(255, 256):
for pointer in guesses:
    guess = prefix + pointer
    p = process(["./eVMoji", "code.bin"])
    p.recvuntil("\xf0\x9f\x8f\xb3\xef\xb8\x8f\n") # white flag emoji
    p.send(guess)
    out = p.recvall() # we have the leaked value in the output
    # if len(out) > 0 and out != b'\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00': 
    # input_bits, result_bits =
    print(out)
    # print(binascii.hexlify(pointer).decode('utf-8'), process_output(out))
  # print(int.from_bytes(pointer, "little"), int.from_bytes(pointer, "big"))
  # print(int.from_bytes(target, "little"), int.from_bytes(target, "big"))
