#!/usr/bin/env python3.8
# Note: The server is running a version of pyopenssl patched with this:
# https://github.com/pyca/pyopenssl/pull/897
# Attacking the pyopenssl wrapper code is not the intended solution.
import OpenSSL.crypto as crypto

welcome = '''=============== WELCOME TO THE RSA TEST SERVICE ===============
You can send me a message and a key to decrypt it!
If your setup works correctly, you will receive a flag as a reward!
But wait, it is quite noisy here!
'''
question_to_ask = b"Hello! Can you give me the flag, please? I would really appreciate it!"


print(welcome)
print("Please give me your private key in PEM format:")
key = ""
while x := input():
    key += x + "\n"

message = input("Now give me your message: ")
message = b"Quack! Quack!"
print("Did you say '" + message.decode() + "'? I can't really understand you, the ducks are too loud!")


key = crypto.load_privatekey(crypto.FILETYPE_PEM, key)
assert key.check()
numbers = key.to_cryptography_key().private_numbers()

d = numbers.d
N = numbers.p * numbers.q

if pow(int.from_bytes(message, "big"), d, N) == int.from_bytes(question_to_ask, "big"):
    print("CSCG{DUMMY_FLAG}")
else:
    print("That was not kind enough!")
