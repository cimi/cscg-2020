package main

import (
	"fmt"
	"os"
	"sync"
)

// var dst uint32 = 0x5e840ef4 // little endian
// var dst uint32 = 0xf40e845e // big endian
// var dst = 0x7a21702f // what the computer tells me to use old
// var dst = 0x85de8fd0 // what the computer tells me to use new
var dst = 0xf40e845e // shrug emoji
// var mix uint32 = 0x2083b8ed // little endian
// var mix uint32 = 0xedb88320 // big endian
// var mix = 0x04c11db7 // what the computer tells me to use
var mix = 0xedb88320
var acc = 0xffffffff // same in both

var size = 32
var chunks = 16
var checked = 0

var dstStr = "11110100000011101000010001011110" // leaked dst
// var dstStr = "01111010001000010111000000101111" // reversed dst

var mixStr = "11101101101110001000001100100000" // leaked mixer
// var mixStr = "00000100110000010001110110110111" // reversed mixer

func main() {
	bruteForce()
}

func bruteForce() {
	var wg sync.WaitGroup
	for i := 0; i < chunks; i++ {
		start, end := makeRange(i)
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			checkRange(start, end)
		}(start, end)
	}
	wg.Wait()
	fmt.Printf("Checked %d numbers\n", checked)
	os.Exit(1)
}

// print the value after going through the mixer to create a reference table
// and compare against the output of this program to find the difference

func makeRange(i int) (start, end int) {
	size := 4294967296 / chunks
	return size * i, size * (i + 1)
}

func checkRange(start, end int) {
	for i := start; i < end; i++ {

		if mixer(i) == dst {
			fmt.Printf("\n🎉🎉 Solution: 0x%08x 0x%08x %s\n\n", i, bitsToNumReverse(numToBits(i)), numToBits(i))
			// os.Exit(0)
		}
	}
	fmt.Printf("⭕️ 0x%08x - 0x%08x (%10d - %10d) \n", start, end-1, start, end-1)
	checked += int(end - start)
}

func mixer(num int) int {
	// everything is little endian, verified by comparing bits
	// we need to make sure we compare against the little endian (in bits)
	// representation of the target numbers, then we can convert back?
	res := acc
	for b := size - 1; b >= 0; b-- {
		// fmt.Printf("%2d %s %s\n", b, numToBits(res), numToBits(mix))
		if ((num >> b) & 1) == (res & 1) {
			res = res >> 1
		} else {
			res = (res >> 1) ^ mix
		}
	}
	return res
}

// func mixerRight(num int) int {
// 	// we constrain to uint32 so we can safely shift right
// 	res := uint32(acc)

// 	for b := 0; b < size; b++ {
// 		if ((num >> b) & 1) == (res & 1) {
// 			res = res >> 1
// 		} else {
// 			res = (res >> 1) ^ mix
// 		}
// 		num = num >> 1
// 	}
// 	return res
// }

/*
dst = 0x5e840ef4 # *136
# >>> '{0:032b}'.format(0x5e840ef4)
# '01011110100001000000111011110100'
# '01011110111000000100001011110100'

# '00001011111100010111101110100001'

acc = 0xffffffff # *140
# '11111111111111111111111111111111'

mix = 0x2083b8ed # *128
# >>> '{0:032b}'.format(0x2083b8ed)
# '00100000100000111011100011101101'

>>> bin(int.from_bytes(b"\x20\x83\xb8\xed", "little"))
'0b11101101101110001000001100100000'

>>> bin(int.from_bytes(b"\x5e\x84\x0e\xf4", "little"))
'0b11110100000011101000010001011110'
we need to use the big endian representation of this number, which is
\xed\xb8\x83
0xedb88320

# '00100000100000111011100011101101'
# '11111111111111111111111111111111'
# '01011110100001000000111011110100'
*/
