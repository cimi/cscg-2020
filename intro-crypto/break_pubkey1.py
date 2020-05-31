from Crypto.Util.number import long_to_bytes
from Crypto.PublicKey import RSA

public_key = RSA.importKey(open('files/crypto1.pem', 'r').read())

print(public_key.n)
print(public_key.e)

message = open('files/message1.txt', 'r').read()

# we need to find d so we can decrypt the message
