import base64
import hashlib
from Crypto import Random
from Crypto.Cipher import AES
from Crypto.Protocol.KDF import PBKDF2
from Crypto.Hash import SHA1
from Crypto.Random import get_random_bytes

# monodis --output=decompiled-reme.txt ReMe.dll

ciphertext = "D/T9XRgUcKDjgXEldEzeEsVjIcqUTl7047pPaw7DZ9I="

password = b'A_Wise_Man_Once_Told_Me_Obfuscation_Is_Useless_Anyway'
salt = b'Ivan Medvedev'
keys = PBKDF2(password, salt, 64, count=1000, hmac_hash_module=SHA1)
key = keys[:32]
iv = keys[32:48]

# Don't work: AES.MODE_CCM, AES.MODE_CTR, AES.MODE_ECB, AES.MODE_OCB, AES.MODE_SIV
# Tried: AES.MODE_CFB,  AES.MODE_EAX, AES.MODE_GCM, AES.MODE_OFB, AES.MODE_OPENPGP
for mode in [AES.MODE_CBC] :
  aes = AES.new(key, mode, iv)
  decrypted = aes.decrypt(base64.b64decode(ciphertext))
  print(decrypted.decode("utf-8"))
