from pwn import *
import OpenSSL.crypto as crypto
import sys

m = int.from_bytes(b"Hello! Can you give me the flag, please? I would really appreciate it!", "big")
c = int.from_bytes(b"Quack! Quack!", "big")

def is_valid_decryption(key, c, m):
  numbers = key.to_cryptography_key().private_numbers()
  d = numbers.d
  N = numbers.p * numbers.q
  return pow(c, d, N) == m

with open("crafty.pem", "r") as kf:
  key_pem = kf.read()
  key = crypto.load_privatekey(crypto.FILETYPE_PEM, key_pem)
  if not key.check():
    print("Key does not pass OpenSSL validation!")
    sys.exit(1)
  if not is_valid_decryption(key, c, m):
    print("Key does not decrypt the cipher correctly!")
    sys.exit(1)
  print("Great success! Key works, now getting the flag...")

  p = remote("hax1.allesctf.net", 9400)
  print(p.recvuntil("Please give me your private key in PEM format:").decode("utf-8"))
  p.send(key_pem + "\n")  
  print(p.recvuntil("Now give me your message:").decode("utf-8"))
  p.send("Quack? Quack?\n")
  print(p.recvall().decode("utf-8"))
  # CSCG{下一家烤鴨店在哪裡？}
