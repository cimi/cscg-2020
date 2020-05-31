
The server asks for a private key and an encrypted payload. It then overwrites the encrypted payload we submit with a hardcoded string ("Quack! Quack!"), then uses the key we gave it to decrypt it and compares it to a fixed target string.

```python
key = crypto.load_privatekey(crypto.FILETYPE_PEM, key)
assert key.check()
numbers = key.to_cryptography_key().private_numbers()

d = numbers.d
N = numbers.p * numbers.q

if pow(int.from_bytes(message, "big"), d, N) == int.from_bytes(question_to_ask, "big"):
    print("CSCG{DUMMY_FLAG}")
else:
    print("That was not kind enough!")
```

The key is checked with openssl and the challenge says we should not focus on attacking the library.
We need to construct a valid key that decrypts the fake cipher to the target message.
The key doesn't need to be secure and we control all parameters - it just needs to pass openssl validation.

Validation checks that p and q are prime and independently that `d*e` is congruent to 1 `| mod (p-1)*(q-1)`.
Since we don't encrypt anything (the ciphertext is fixed), e is not used so we can pick any value that passes validation.

Essentially we have two independent pairs of params (p,q) and (d,e); the pairs have the congruence relationship:
`e * d = i*(p-1)*(q-1) + 1`, where i can be any natural number.

p and q need to be prime; e and d don't have to be prime, they just have to be natural numbers.

```python
# pick small integer values for d such that kN has known factors
def find_d_and_factored_kN(c,m):
  for d in range(3, 100):
    kn = pow(c, d) - m 

    f = FactorDB(kn)
    f.connect()
    if f.get_status() == 'FF':
      print(d, " Win!")  
      print(kn)
      print(f.get_factor_list())
    else:
      print(f.get_status())
```

We filter for numbers that have known factorings in factordb. We find two values that work for d (14 and 97) along with the associated kN and its factors, so now we have p and q.

We need to compute `e` so the key passes openssl validation.

```python
def find_e(d, p, q):
  # pick k such that (k*(p-1)*(q-1) + 1) / d is an integer
  found = False
  for k in range(3, 100):
    if (k * (p-1) * (q-1) + 1) % d == 0:
      print(k, "Win!")
      found = True
      e = (k * (p-1) * (q-1) + 1) // d
      return e
  if not found:
    print(d, "No win!")
    return
```

Finally, putting all this together:

```python
def try_genkey(p, q, d):  
  e = find_e(p, q, d)
  if e == None:
    print("Invalid parameters, no e found!")
    return
  
  with open("crafty.pem", "wb") as kf:
    key = RSA.construct((p*q, e, d, p, q), False)
    kf.write(key.exportKey())
    print("Key saved as crafty.pem")
```

We can check the key decrypts the message correctly:

```python
def is_valid_decryption(key, c, m):
  numbers = key.to_cryptography_key().private_numbers()
  d = numbers.d
  N = numbers.p * numbers.q
  return pow(c, d, N) == m
```

Now we use pwntools to do the server i/o and get the flag:

```python
p = remote("hax1.allesctf.net", 9400)
print(p.recvuntil("Please give me your private key in PEM format:").decode("utf-8"))
p.send(key_pem + "\n")  
print(p.recvuntil("Now give me your message:").decode("utf-8"))
p.send("Quack? Quack?\n")
print(p.recvall().decode("utf-8"))
```

```console
root@14a2c8b9938e:/pwd/rsa-service# python3 solve.py
Great success! Key works, now getting the flag...
[+] Opening connection to hax1.allesctf.net on port 9400: Done
=============== WELCOME TO THE RSA TEST SERVICE ===============
You can send me a message and a key to decrypt it!
If your setup works correctly, you will receive a flag as a reward!
But wait, it is quite noisy here!

Please give me your private key in PEM format:

Now give me your message:
[+] Receiving all data: Done (122B)
[*] Closed connection to hax1.allesctf.net port 9400
 Did you say 'Quack! Quack!'? I can't really understand you, the ducks are too loud!
CSCG{下一家烤鴨店在哪裡？}
```