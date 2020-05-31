#from pwn import *
from tqdm import tqdm

dst = 0x5e840ef4 # *136
# >>> '{0:032b}'.format(0x5e840ef4)
# '01011110100001000000111011110100'

acc = 0xffffffff # *140
# '11111111111111111111111111111111'

mix = 0x2083b8ed # *128
# >>> '{0:032b}'.format(0x2083b8ed)
# '00100000100000111011100011101101'

# '00100000100000111011100011101101'
# '11111111111111111111111111111111'
# '01011110100001000000111011110100'

def bits(n):
  while n:
    b = n & (~n+1)
    yield b
    n ^= b

for x in tqdm(range(0, 0xffffffff)):
  # if x%1024*1024 == 0:
  #   print(x, int.to_bytes(x, 8, "little"))
  acc = 0xffffffff # *140
  res = ""
  for b in range(0,32):
    is_set = (x & (1 << (31 - b))) > 0
    if is_set:
      cur_bit = 1
    else:
      cur_bit = 0
    if cur_bit == (acc & 1):
      acc = acc >> 1
    else:
      acc = (acc >> 1)^mix
  if acc == dst:
    print("WIN!", acc, dst, x)




# for each bit in source:
#   if bit == last_bit(acc): # this keeps flipping with every xor, initially last_bit(acc) is 1
#     acc = acc >> 1
#   else:
#     acc = (acc >> 1)^mix
# assert acc == dst

# if mixer[i] != target[i]
# each bit different from mixer in target must have an odd number of bits set in remainder of source

# pick number of possible valid flips for each position
# for some of them we must have a restriction to be odd

# first bit will always xor against zero (mixer is zero)
# because the accumulator is shifted before xor, regardless of the input first bit the result on the accumulator on the first bit is zero

# second bit needs to flip from zero in mixer to one in dst
# we don't know if we've xored once before and it doesn't affect this one since it would xor zero and one and yield one anyway

# we need to xor in the second one because otherwise it will be set to zero and mixer won't recover it
# second bit is one! - or not - it can also flip the last bit so we don't know when we mix

# third bit needs to stay zero in target
# this will keep flipping on every xor because mixer is 1

# for each bit that's zero in the mixer and one in the target we can't xor on that position
# for each bit that's zero in the target and one in the mixer, we need to xo and odd number of times after its position
  # because it flips between zero and one every time it xors
