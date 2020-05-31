set -e
mkdir -p spectrograms/parts
for wav in parts/*.wav; do
  # hexdump -C $wav > $wav.hex
  sox $wav -n spectrogram -o spectrograms/$wav.png
  # stegolsb wavsteg -r -i $wav -o $wav.01.bin -n 1 -b 10000 &> /dev/null 
  # stegolsb wavsteg -r -i $wav -o $wav.02.bin -n 2 -b 20000 &> /dev/null 
  # stegolsb wavsteg -r -i $wav -o $wav.03.bin -n 3 -b 30000 &> /dev/null 
  # stegolsb wavsteg -r -i $wav -o $wav.04.bin -n 4 -b 40000 &> /dev/null 
  # hexdump -C $wav.01.bin > $wav.01.bin.hex
  # hexdump -C $wav.02.bin > $wav.02.bin.hex
  # hexdump -C $wav.03.bin > $wav.03.bin.hex
  # hexdump -C $wav.04.bin > $wav.04.bin.hex
  echo "✅ $wav"
done

# xor parts/matrix-01.wav.01.bin parts/matrix-02.wav.01.bin > parts/matrix-xor-01-02.01.bin
# xor parts/matrix-01.wav.02.bin parts/matrix-02.wav.02.bin > parts/matrix-xor-01-02.02.bin
# xor parts/matrix-01.wav.03.bin parts/matrix-02.wav.03.bin > parts/matrix-xor-01-02.03.bin
# xor parts/matrix-01.wav.04.bin parts/matrix-02.wav.04.bin > parts/matrix-xor-01-02.04.bin

# xor parts/matrix-02.wav.01.bin parts/matrix-03.wav.01.bin > parts/matrix-xor-02-03.01.bin
# xor parts/matrix-02.wav.02.bin parts/matrix-03.wav.02.bin > parts/matrix-xor-02-03.02.bin
# xor parts/matrix-02.wav.03.bin parts/matrix-03.wav.03.bin > parts/matrix-xor-02-03.03.bin
# xor parts/matrix-02.wav.04.bin parts/matrix-03.wav.04.bin > parts/matrix-xor-02-03.04.bin

# hexdump -C parts/matrix-xor-01-02.01.bin > parts/matrix-xor-01-02.01.bin.hex
# hexdump -C parts/matrix-xor-01-02.02.bin > parts/matrix-xor-01-02.02.bin.hex
# hexdump -C parts/matrix-xor-01-02.03.bin > parts/matrix-xor-01-02.03.bin.hex
# hexdump -C parts/matrix-xor-01-02.04.bin > parts/matrix-xor-01-02.04.bin.hex

# hexdump -C parts/matrix-xor-02-03.01.bin > parts/matrix-xor-02-03.01.bin.hex
# hexdump -C parts/matrix-xor-02-03.02.bin > parts/matrix-xor-02-03.02.bin.hex
# hexdump -C parts/matrix-xor-02-03.03.bin > parts/matrix-xor-02-03.03.bin.hex
# hexdump -C parts/matrix-xor-02-03.04.bin > parts/matrix-xor-02-03.04.bin.hex

# Th3-R3D-P1ll?




# 0110 1110 
# 0010 0001 
# 0100 0011 
# 0011 0011
# 0101 1111 
# 0101 0000
# 0101 0111 
# 0011 1111

# 01101110 
# 00100001 
# 01000011 
# 00110011
# 01011111 
# 01010000
# 01010111 
# 00111111

bs = ["01101110","00100001","01000011","00110011","01011111","01010000","01010111","00111111"]

# 0x6e 0x21
# 0x43 0x33
# 0x5f 0xa0


#!1001 0001 1100 1100
n!C3_PW?

63?
11100110
# ae 21 
