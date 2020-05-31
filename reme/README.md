## First flag - AES decryption

I used dotPeek from JetBrains to decompile the dll but probably dnSpy, monodis or others work just as well.

After decompiling, we notice the first argument is compared to the decryption of a cipher string - that's the first flag: 

```cs
    if (args.Length == 0) {
        Console.WriteLine("Usage: ReMe.exe [password] [flag]");
        Environment.Exit(-1);
    }
    if (args[0] != StringEncryption.Decrypt("D/T9XRgUcKDjgXEldEzeEsVjIcqUTl7047pPaw7DZ9I=")) {
        Console.WriteLine("Nope");
        Environment.Exit(-1);
    } else {
        Console.WriteLine("There you go. Thats the first of the two flags! CSCG{{{0}}}", (object) args[0]);
    }
```

Decrypt shows us we're looking at a base64 encoded string of AES ciphertext and we're given the password and the salt, so we can write our own code to decrypt the string.

```cs
    public static string Decrypt(string cipherText)
    {
      string password = "A_Wise_Man_Once_Told_Me_Obfuscation_Is_Useless_Anyway";
      cipherText = cipherText.Replace(" ", "+");
      byte[] numArray = Convert.FromBase64String(cipherText);
      using (Aes aes = Aes.Create())
      {
        Rfc2898DeriveBytes rfc2898DeriveBytes = new Rfc2898DeriveBytes(password, new byte[13]
        {
          (byte) 73,
          (byte) 118,
          (byte) 97,
          (byte) 110,
          (byte) 32,
          (byte) 77,
          (byte) 101,
          (byte) 100,
          (byte) 118,
          (byte) 101,
          (byte) 100,
          (byte) 101,
          (byte) 118
        });
        aes.Key = rfc2898DeriveBytes.GetBytes(32);
        aes.IV = rfc2898DeriveBytes.GetBytes(16);
        using (MemoryStream memoryStream = new MemoryStream())
        {
          using (CryptoStream cryptoStream = new CryptoStream((Stream) memoryStream, aes.CreateDecryptor(), CryptoStreamMode.Write))
          {
            ((Stream) cryptoStream).Write(numArray, 0, numArray.Length);
            ((Stream) cryptoStream).Close();
          }
          cipherText = Encoding.get_Unicode().GetString(memoryStream.ToArray());
        }
      }
      return cipherText;
    }
  }
```

Here's python code to decrypt this string and get the first flag:

```python
import base64
import hashlib
from Crypto import Random
from Crypto.Cipher import AES
from Crypto.Protocol.KDF import PBKDF2
from Crypto.Hash import SHA1
from Crypto.Random import get_random_bytes

ciphertext = "D/T9XRgUcKDjgXEldEzeEsVjIcqUTl7047pPaw7DZ9I="

password = b'A_Wise_Man_Once_Told_Me_Obfuscation_Is_Useless_Anyway'
salt = b'Ivan Medvedev'
keys = PBKDF2(password, salt, 64, count=1000, hmac_hash_module=SHA1)
key = keys[:32]
iv = keys[32:48]

for mode in [AES.MODE_CBC] :
  aes = AES.new(key, mode, iv)
  decrypted = aes.decrypt(base64.b64decode(ciphertext))
  print(decrypted)
```

```console
root@14a2c8b9938e:/pwd/reme# python3 ./solve_reme1.py
CanIHazFlag? 
=> CSCG{CanIHazFlag?}
```

## Second flag - dynamic code loading and AES

```cs
 private static void Main(string[] args)
    {
      Program.InitialCheck(args);
			// get method bytecode
      byte[] ilAsByteArray = ((MethodBase) typeof (Program).GetMethod("InitialCheck", (BindingFlags) 40)).GetMethodBody().GetILAsByteArray();
			// load dll in memory stream
      byte[] self = File.ReadAllBytes(Assembly.GetExecutingAssembly().get_Location());
			// find string in dll?
      int[] numArray = self.Locate(Encoding.get_ASCII().GetBytes("THIS_IS_CSCG_NOT_A_MALWARE!"));
      MemoryStream memoryStream = new MemoryStream(self);
			// jump to right after string
      ((Stream) memoryStream).Seek((long) (numArray[0] + Encoding.get_ASCII().GetBytes("THIS_IS_CSCG_NOT_A_MALWARE!").Length), (SeekOrigin) 0);
			// take the rest of the file
      byte[] bytesToBeDecrypted = new byte[((Stream) memoryStream).get_Length() - ((Stream) memoryStream).get_Position()];
      ((Stream) memoryStream).Read(bytesToBeDecrypted, 0, bytesToBeDecrypted.Length);
			// Decrypt the bytes and execute them?! We call the check method on whatever we decrypt
      ((MethodBase) Assembly.Load(Program.AES_Decrypt(bytesToBeDecrypted, ilAsByteArray)).GetTypes()[0].GetMethod("Check", (BindingFlags) 24)).Invoke((object) null, new object[1]
      {
        (object) args
      });
    }
```

We can decrypt using a python script, similar to what we did in the previous step (see solve_reme2.py). We extract the inner dll and decompile it, we see that the second flag is 

```cs
bool flag3 = "CSCG{" + array[0] == "CSCG{n0w" 
    && array[1] == "u" 
    && array[2] == "know" 
    && array[3] == "st4t1c" 
    && array[4] == "and" 
    && Inner.CalculateMD5Hash(array[5]).ToLower() == "b72f3bd391ba731a35708bfd8cd8a68f" 
    && array[6] == "dotNet" 
    && array[7] + "}" == "R3333}";
    if (flag3) {
        Console.WriteLine("Good job :)");
    }
```

We can lookup the hash in crackstation, the input is `dynamic`. So we have the second flag!

```console
C:\Users\alexc\Downloads> dotnet ReMe.dll CanIHazFlag? n0w_u_know_st4t1c_and_dynamic_dotNet_R3333                    
There you go. Thats the first of the two flags! CSCG{CanIHazFlag?}
Good job :)

=> CSCG{n0w_u_know_st4t1c_and_dynamic_dotNet_R3333}
```

