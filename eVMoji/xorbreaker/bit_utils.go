package main

import "fmt"

func printBits(num int, expected string) {
	fmt.Printf("0x%08x -> ", num)
	bits := ""
	for b := 0; b < size; b++ {
		bitval := (num >> b) & 1
		bits += fmt.Sprintf("%d", bitval)
	}
	if bits == expected {
		fmt.Println(bits)
		return
	}
	fmt.Printf("0x%08x != 0x%08x (%s != %s)\n", num, bitsToNum(expected), bits, expected)
}

func numToBits(num int) string {
	// also print in reverse order
	bits := ""
	for b := 0; b < size; b++ {
		bitval := (num >> b) & 1
		bits += fmt.Sprintf("%d", bitval)
	}
	return bits
}

func bitsToNum(s string) int {
	num := 0
	for idx := range s {
		if s[idx] == byte('1') {
			num += pow(2, idx)
		}
	}
	return num
}

func bitsToNumReverse(s string) int {
	num := 0
	for idx := range s {
		if s[len(s)-1-idx] == byte('1') {
			num += pow(2, idx)
		}
	}
	return num
}

func pow(a, b int) int {
	res := 1
	for i := 1; i <= b; i++ {
		res = res * a
	}
	return res
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
