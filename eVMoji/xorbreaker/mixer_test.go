package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

var testCases = make(map[string]string)

func init() {
	testfile, err := os.Open("mixer-test-cases.txt")
	defer testfile.Close()
	check(err)
	reader := bufio.NewReader(testfile)

	var line string
	for {
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		check(err)
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		s1 := strings.Split(line, " ")[1]
		s1 = strings.TrimRight(s1, "| \n")
		binstrs := strings.Split(s1, "|")
		if len(binstrs) != 2 {
			panic("Failure reading test cases")
		}
		testCases[binstrs[0]] = binstrs[1]
		// fmt.Printf(" > %s\n", s1)
		// for _, binstr := range strings.Split(s1, "|") {
		// 	fmt.Printf(" > %s\n", binstr)
		// }
	}
}

// 0001 0000 
// 0000 0000 0000 0000 0000 0001 0000 0000

// fffbffff 00000000011111111111111111111111|00000000001111111111111111111111|11101101101001110111110011011111|
// 11111111111111111111110111111111
// position 11 start flipping
// in mine, position 23 starts flipping
// what happens at position 11 in mine?

func TestMixer(t *testing.T) {
	fmt.Printf("Using mixer %s\n", numToBits(mix))
	for bitsIn, bitsOut := range testCases {
		in := bitsToNum(bitsIn)
		fmt.Printf("Checking %s as 0x%08x: ", bitsIn, in)

		expected := bitsToNumReverse(bitsOut)
		actual := mixer(in)
		if actual != expected {
			fmt.Println("⚠️")
			t.Errorf("Failed 0x%08x: 0x%08x != 0x%08x\n%s\n%s\n", in, actual, expected, numToBits(actual), numToBits(expected))
		} else {
			fmt.Println("✅")
		}
	}
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func TestBitFunctions(t *testing.T) {
	num := 0xdeadbeef
	if num != bitsToNum(numToBits(num)) {
		t.Errorf("Bit conversion failed!")
	}
	if bitsToNumReverse(numToBits(num)) != bitsToNum(reverse(numToBits(num))) {
		t.Errorf("Bit reversal failed: 0x%08x != 0x%08x!", bitsToNumReverse(numToBits(num)), bitsToNum(reverse(numToBits(num))))
	}

	leakedMixerBits := "11101101101110001000001100100000"
	leakedDstBits := "11110100000011101000010001011110"
	if reverse(leakedMixerBits) != numToBits(mix) {
		t.Errorf("Mixer value 0x%08x is %s != %s", mix, leakedMixerBits, numToBits(mix))
	}

	if reverse(leakedDstBits) != numToBits(dst) {
		t.Errorf("Target value 0x%08x is %s != %s", dst, leakedDstBits, numToBits(dst))
	}
}

func TestSolution(t *testing.T) {
	solutionBits := "00111111011011000011000001101100"
	if mixer(bitsToNum(solutionBits)) != dst {
		t.Errorf("Wrong solution")
	}
}
