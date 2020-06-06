package main

import (
	"fmt"
	"os"
	"sync"
)

var size = 32
var chunks = 16
var checked = 0

var dstStr = "11110100000011101000010001011110" // leaked dst
var dst = 0xf40e845e
var mixStr = "11101101101110001000001100100000" // leaked mixer
var mix = 0xedb88320

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
	res := 0xffffffff
	for b := size - 1; b >= 0; b-- {
		if ((num >> b) & 1) == (res & 1) {
			res = res >> 1
		} else {
			res = (res >> 1) ^ mix
		}
	}
	return res
}
