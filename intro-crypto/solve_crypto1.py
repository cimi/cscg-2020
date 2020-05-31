from fractions import gcd
from Crypto.Util.number import long_to_bytes
from Crypto.PublicKey import RSA

public_key = RSA.importKey(open('files/crypto1.pem', 'r').read())
N = public_key.n
e = public_key.e

for i in range(2, 1000000):
  if N % i == 0:
    p = i
    print("Factor found: " + str(i))
    break

q = N / p
phin = (p-1) * (q-1)

def xgcd(a, b):
  prevx, x = 1, 0; prevy, y = 0, 1
  while b:
    q = a / b
    x, prevx = prevx - q*x, x
    y, prevy = prevy - q*y, y
    a, b = b, a % b
  return a, prevx, prevy

def modinv(a, m):
    """return x such that (x * a) % b == 1"""
    g, x, _ = xgcd(a, m)
    if g != 1:
        raise Exception('gcd(a, b) != 1')
    return x % m

d = modinv(e, phin)

m = int(open('files/message1.txt', 'r').read())
decrypted = pow(m, d, N)
print(long_to_bytes(decrypted))
