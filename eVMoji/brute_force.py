# -*- coding: UTF-8 -*-
from pwn import *
from tqdm import trange, tqdm
import string

nums = []

# we used to use only the numbers available in the source
with open("analysis/numbers.txt", "r") as nf:
  lines = nf.readlines()
  for l in lines:
    if int(l) < 128:
      nums.append(int(l)) 

def generate_program_text_start(nums):
  with open("tmp/generated.txt", "w") as pf:
    pf.write("I0\n")
    pf.write("I2R\n")
    pf.write("I" + nums[0] + "I1W\n")
    pf.write("I" + nums[1] + "I1W\n")
    pf.write("I" + nums[0] + "M0S\n")
    pf.write("I" + nums[1] + "SC\n")
    pf.write("I1WD")

def generate_program_text_cyclic(nums):
  with open("tmp/generated.txt", "w") as pf:
    pf.write("I0\n")
    pf.write("I1R\n")
    pf.write("I0\n")
    pf.write(">1!+*0>" + nums[0] + "+T236>1*128SI0\n")
    pf.write("I0\n")
    pf.write("I1WD")

# for each pair generate the machine code text
# then use the go program to convert the text to binary
# then use the VM to execute the code and fish out the results
# finally, present the results :)
def brute_force():
  with open("analysis/pairs2.txt", "r") as addrfile:
    addrlines = addrfile.readlines()
    for al in addrlines:
      nums = al.strip().split(",")
      # generate text instructions 
      generate_program_text_start(nums)

      # build program binary
      os.system("./emoji build tmp/generated.txt tmp/generated.bin")
      
      # now try everything and store keys that match
      for i in range(1, 127):
        p = process(["./eVMoji", "tmp/generated.bin"])
        guess = chr(i)
        p.send(guess + "\n")
        out = p.recvall()
        if len(out) >= 3 and out[2] == i:
          if len(out) == 3:
            print(nums, guess, out[0], out[1], out[2])
          else:
            print(out)

some_emoji = [emoji.encode('utf-8') for emoji in [chr(i) for i in range(0x1F601, 0x1F64F)]]
prefix = b"n3w_ag3_v1rtu4liz4t1on_"

def brute_force_unfeasible():  
  for x in trange(255, desc="first"):
    for y in trange(255, desc="second"):
      for z in trange(255, desc="third"):
        for t in trange(255, desc="fourth", leave=False):
          p = process(["./eVMoji", "code.bin"])
          guess = prefix + bytes([x, y, z, t])
          p.recvuntil("\xf0\x9f\x8f\xb3\xef\xb8\x8f\n") # white flag emoji
          p.send(guess + b"\n")
          out = p.recvall()
          if out != b"Gotta go cyclic \xe2\x99\xbb\xef\xb8\x8f\n":
            print(guess, out)

# suffix = bytes([0xE2, 0x99, 0xBB]
suffix = b"abc"
def brute_force_cyclic():
  for i in range(0, 31):
    matches = str(i) + ": "
    generate_program_text_cyclic([str(i)])
    os.system("./emoji build tmp/generated.txt tmp/generated.bin")
    for x in range(0, 127):
      p = process(["./eVMoji", "tmp/generated.bin"])
      guess = chr(x)
      p.send(guess + "\n")
      out = p.recvall()
      if len(out) > 0:
        # matches += out.decode('utf-8').strip()
        matches += str(ord(out.decode('utf-8'))) + ", "
    print(matches)


def brute_force_alpha():
  charset = string.ascii_letters + string.digits + string.punctuation
  for x in charset:
    for y in charset:
      print("Guessing", x, y)
      for z in charset:
        for t in charset:
          p = process(["./eVMoji", "code.bin"])
          guess = "n3w_ag3_v1rtu4liz4t1on_"+x+y+z+t
          p.recvuntil("\xf0\x9f\x8f\xb3\xef\xb8\x8f\n") # white flag emoji
          p.send(guess + "\n")
          out = p.recvall()
          if out != b"Gotta go cyclic \xe2\x99\xbb\xef\xb8\x8f\n":
            print(guess, out)

def generate_program_text_star(nums):
  with open("tmp/generated.txt", "w") as pf:
    # pf.write("I"+nums[0]+"I"+nums[2]+"R\n")
    # pf.write("I0\n")
    # pf.write("*"+nums[1]+"\n")  
    # pf.write("I20WD")
    pf.write("I0")
    pf.write("I27")
    pf.write("R")
    pf.write("I0I242M0SI156SCI234M1SI217SCI130M2SI245SCI54M3SI105SCI142M4SI239SCI18M5SI117SCI24M6SI43SCI115M7SI44SCI123M8SI13SCI17M9SI32SCI91M10SI41SCI105M11SI29SC")
    pf.write("I56M12SI77SCI138M13SI190SCI176M14SI220SCI139M15SI226SCI142M16SI244SCI131M17SI183SCI246M18SI130SCI196M19SI245SCI57M20SI86SCI245M21SI155SCI162M22SI253SC")
    pf.write("I0T102I187I25WD")
    pf.write("*140>8")
    pf.write("I"+nums[0])
    pf.write("T102")
    pf.write("I187I25WD")
    pf.write("I0I27WD")

def check_pointer():
  guess = "n3w_ag3_v1rtu4liz4t1on_BBBB"
  for x in range(0, 256):
    if x != 203:
      generate_program_text_star([str(x)])
    else:
      generate_program_text_star([str(x*2)+">1"])
    os.system("./emoji build tmp/generated.txt tmp/generated.bin")
    p = process(["./eVMoji", "tmp/generated.bin"])
    p.send(guess)
    out = p.recvall()
    if out != b'tRy hArder! \xf0\x9f\x92\x80\xf0\x9f\x92\x80\xf0\x9f\x92\x80\n':
      print(x, out)
    elif x % 16 == 0:
      print(x)

check_pointer()
# n3w_ag3_v1rtu4liz4t1on_xxxx


def brute_force_pointer():
  with open("solve-output.txt", "w") as outfile:
    # x - position, y - pointer, z - length
    z = 5
    for x in trange(32, leave=False):
      for y in trange(32, leave=False):
        for z in range(0, 9):
          generate_program_text_star([str(x), str(y), str(z)])
          os.system("./emoji build tmp/generated.txt tmp/generated.bin")
          p = process(["./eVMoji", "tmp/generated.bin"])
          guess = "\xeb" * z + "\n"
          # guess += "\xec" * 5    
          p.send(guess + "\n")
          out = p.recvall()
          outfile.write(f"pos: {x:2d} || len: {z} || *: {y:2d} => ")
          if out == b'Thats the flag: CSCG':
            outfile.write("JUMP\n")
          elif out == b"\x00" * 20:
            outfile.write("0x00\n")
          elif len(out) == 0:
            outfile.write("NULL\n")
          else:
            outfile.write(f"{out[:20]}\n")

def ok_this_is_ridiculous():
  guess = "n3w_ag3_v1rtu4liz4t1on_BBBB"
  os.system("./emoji build tmp/generated.txt tmp/generated.bin")
  p = process(["./eVMoji", "tmp/generated.bin"])
  p.recvuntil("\xf0\x9f\x8f\xb3\xef\xb8\x8f\n") # white flag emoji
  p.send(guess + "\n")
  out = p.recvall()
  print(out)

# when *i <= len - 3 => '';
# when *i == len - 2 => \x00
# when *i == len - 1 => JUMP!
# Ix*i => 0 if x-y > 4 OR i >= len(input); 

# I6*0, I6*2 => 0
# I6*3, I6*

# 100: caaaUnknown opcode: 8092sh returned exit code 255
# 128: 
# 102: caaash returned exit code 255 
# aaaaUnknown opcode: 30Unknown opcode: 8fb8efUnknown opcode: a383e2Unknown opcode: 30Unknown opcode:


# S - XOR(pop1(), pop2())
# >[i] - shift(pop())
# M[i] - push(chr(i))
# I[i] - push(i)
# *[i] - reads four? bytes then loads the contents of address i on the stack; address seems to be little endian
# T[i] - x = pop(); y = pop() if x == y then move i+2 bytes forward
# ! - x = pop(); push(x); push(x)
# + - push(pop() & 1) 
# C - bitwise or |
# R - read(stdin, dstaddr=pop2(), len=pop1())
# W - write(stdout, srcaddr=pop2(), len=pop1() 

# * 

# generate all numbers
# overwrite data segment with identifiable pattern - cyclic!

# >>> pairs = 
# ...[(242, 156, 'm', 'n'),
# ... (234, 217, '2', '3'),
# ... (130, 245, 'v', 'w'),
# ... (54 , 105, '^', '_'),
# ... (142, 239),
# ... (18 , 117),
# ... (24 ,  43),
# ... (115,  44),
# ... (123,  13),
# ... (17 ,  32),
# ... (91 ,  41, 'r', 's')
# ... (105,  29),
# ... (56 ,  77),
# ... (138, 190),
# ... (176, 220),
# ... (139, 226),
# ... (142, 244),
# ... (131, 183),
# ... (246, 130),
# ... (196, 245),
# ... (57 ,  86),
# ... (245, 155),
# ... (162, 253)]


# (242, 'h', 156, 'V')
# (234, '\n', 217, ' ')
# (130, '\xb8', 245, 'f')
# (54, '\x00', 105, '\x00')
# (142, '\xff', 239, 's')
# (18, '\x00', 117, '\x00')
# (24, '\x00', 43, '\x00')
# (115, '\x00', 44, '\x00')
# (123, '\x00', 13, '\x00')
# (17, '\x00', 32, '\x00')
# (91, '\x00', 41, '\x00')
# (105, '\x00', 29, '\x00')
# (56, '\x00', 77, '\x00')
# (138, '\x0e', 190, ' ')
# (176, 'h', 220, ' ')
# (139, '\xf4', 226, 'c')
# (142, '\xff', 244, ' ')
# (131, '\xed', 183, '\xef')
# (246, 'l', 130, '\xb8')
# (196, 'r', 245, 'f')
# (57, '\x00', 86, '\x00')
# (245, 'f', 155, 'e')
# (162, '\xf0', 253, 'C')
