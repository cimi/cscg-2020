import base64
import re
from Crypto.Cipher import AES
from Crypto.Protocol.KDF import PBKDF2
from Crypto.Hash import SHA1

with open("ReMe.dll", "rb") as f:
  dll_bytes = f.read()
  barr = bytearray(dll_bytes)
  print(re.search(b"THIS_IS_CSCG_NOT_A_MALWARE!", barr).span()) # end: 13339
  # found this string by printing out the value of the variable in the check function
  print(re.search("\x00\x00\x0a\x00\x00\\\x28\x09\x00\x00\x06\x13\x04\x11\x04\x2c".encode(), barr).span())

# extract the contents of the dll after the CSCG string and use as chiphertext
with open("ReMe.dll", "rb") as f:
  discard = f.read(13339)
  ciphertext = f.read()

# extract the IL code for IntialCheck and use as the PBKDF2 password
# this changes when we recompile the DLL so we need to run on the unmodified binary
with open("ReMe.dll", "rb") as f:
  discard = f.read(1224)
  il_code = f.read(289)
  f.read()

password = il_code
# the salt is hardcoded in the binary
salt = b'\x01\x02\x03\x04\x05\x06\x07\x08'
keys = PBKDF2(password, salt, 64, count=1000, hmac_hash_module=SHA1)
key = keys[:32]
iv = keys[32:48]

aes = AES.new(key, AES.MODE_CBC, iv)
decrypted = aes.decrypt(ciphertext)
with open("inner.dll", "wb") as f:
  f.write(decrypted)

# C:\Users\alexc\Downloads> dotnet ReMe.dll CanIHazFlag? n0w_u_know_st4t1c_and_dynamic_dotNet_R3333                    
# There you go. Thats the first of the two flags! CSCG{CanIHazFlag?}
# Good job :)
