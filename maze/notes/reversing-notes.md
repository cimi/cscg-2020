
        *(undefined *)(byteArr + 0x20) = (char)len;
        if (*(longlong **)(serverManager + 0x1c0) != (longlong *)0x0) {
          unknown = (**(code **)(**(longlong **)(serverManager + 0x1c0) + 0x180))();
          len = *(uint *)(byteArr + 0x18);
          if (len < 2) goto LAB_1806d5aa1;
          *(undefined *)(byteArr + 0x21) = unknown;
          while( true ) {
            if ((int)*(uint *)(bytes + 0x18) <= (int)idx) break;
            if (*(uint *)(bytes + 0x18) <= idx) goto LAB_1806d5ab0;
            idxPlus2 = (longlong)(int)idx + 2;
            if (len <= (uint)idxPlus2) goto LAB_1806d5abf;
            *(byte *)((idxPlus2 & 0xffffffff) + 0x20 + byteArr) =
                 *(byte *)((longlong)(int)idx + 0x20 + bytes) ^ (byte)rand;
            len = *(uint *)(byteArr + 0x18);
            if (len < 2) goto LAB_1806d5ace;
            intermediary = ((uint)rand & 0xff) + (uint)*(byte *)(byteArr + 0x21);
            rand = (ulonglong)(byte) ((char)intermediary + (char)(intermediary / 0xff));
            idx = idx + 1;
          }
          if (*(longlong *)(serverManager + 0x138) != 0) {
            uVar2 = System.Net.Sockets.UdpClient$$Send
                              (*(longlong *)(serverManager + 0x138),byteArr,(ulonglong)len,0);
            return CONCAT71((int7)((ulonglong)uVar2 >> 8),1);
          }

    -> allocates byte array with len = len(original) + 2
    -> generates two random numbers, one is stored in result[1], the other one might not be stored anywhere
    -> if length is stored in result[0], the random number can be obtained by xor-ing the length against the first position
    -> for every subsequent position, the in memory random is updated, but not the other one
    -> we add the two randoms together, then we do newRand = (sum % 0xff) + (sum / 0xff); so if the sum overflows we reset and add 1
    
    -> the user secret is appended to the emoji number when sending that data (TODO: see if sends bytes, probably)

    So len of emoji packet should be len(secret) + 1 byte for the emoji + 2 bytes for len and random

    Secret: 9D 64 15 A9 AD 96 12 66

    maybe two bytes for the emoji?



    first two bytes are random
    third byte is in[0] ^ first rand

    third byte should be first byte ^ first random -> for emoji this could be always zero
    fourth byte should should be the emoji value ^ new random
    fifth byte should be first secret byte ^ new random
    etc.

    send 01 (0x17 / 23)
    0000   36 78 73 33 43 8a b1 3d 9f 93 9f 65
    0000   30 cf 75 9d ab 8a c6 92 99 cc c8 69


    Data sent is '0x45 + secret + emoji_id':

*(undefined *)(dstArr + 0x20) = 0x45;
System.Buffer$$BlockCopy(*(undefined8 *)(this + 0x178),0,dstArr,1,8,0);
if (9 < *(uint *)(dstArr + 0x18)) {
    *(undefined *)(dstArr + 0x29) = (char)emoji_id;
    ServerManager$$sendData(this,dstArr,0);
    return;
}

So it looks like this: 0x45 yyyy yyyy e (len=10) 
We get 12 bytes sent which is good! two from randoms plus garbled nonsense
we need to reverse the garbled nonsense generator!


