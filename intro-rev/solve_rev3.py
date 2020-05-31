# (0xa + idx) xor char
def transform(num, idx):
  return num ^ (0xa + idx) - 2

def reverse(s):
  result = []
  for idx, c in enumerate(s):
    result.append(chr((ord(c) + 2) ^ (0xa + idx)))
  return "".join(result)
