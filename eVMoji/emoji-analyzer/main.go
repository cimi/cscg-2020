package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var digits = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
var re2 = regexp.MustCompile(`\d+|[^\d]+`)
var encoding = map[string]string{}

func encodeFile(path string) []*Instruction {
	instructions := []*Instruction{}
	decoded, err := ioutil.ReadFile(path)
	check(err)
	lines := strings.Split(string(decoded), "\n")
	for _, l := range lines {
		if l == "" {
			continue
		}
		instructions = append(instructions, instr(c.encodeInstr(l)))
	}
	return instructions
}

func cyclicData() []byte {
	input := "aaaabaaacaaadaaaeaaafaaagaaahaaaiaaajaaakaaalaaamaaanaaaoaaapaaa"
	result := []byte{}
	for i := 0; i < 0x200; i++ {
		if i < len(input) {
			result = append(result, byte(input[i]))
		} else {
			result = append(result, 0)
		}
	}
	return result
}

func orderedData() []byte {
	result := []byte{}
	for i := 0; i < 0x200; i++ {
		if i < 128 {
			result = append(result, byte(i))
		} else {
			result = append(result, 0)
		}
	}
	return result
}

func craftyData() []byte {
	result := []byte{}
	for i := 0; i < 0x200; i++ {
		if i < 48 {
			result = append(result, byte(i % 10 + 48))
		} else if i < 255 {
			result = append(result, byte(((i - 48) % 25) + 65))
		} else {
			result = append(result, 0)
		}
	}
	return result
}

func main() {
	cmd := os.Args[1]
	switch cmd {
	case "analyse":
		name := os.Args[2]
		original := load(name)
		original.writeDecoded(fmt.Sprintf("analysis/%s.txt", name))
		ioutil.WriteFile(fmt.Sprintf("analysis/%s.dump", name), []byte(original.pretty(true)), 0644)
		ioutil.WriteFile(fmt.Sprintf("analysis/%s.emoji", name), []byte(original.pretty(false)), 0644)

		// ioutil.WriteFile("analysis/decoded.txt", []byte(original.decode()), 0644)
	case "build":
		srcfile := os.Args[2]
		dstfile := os.Args[3]
		// we load the original so we get the encoding table
		p := load("code.bin")
		// we overwrite the code with ours and write a new binary
		// p.data = craftyData()
		// p.data = orderedData()
		p.code = encodeFile(srcfile)
		p.writeBinary(dstfile)
	default:
		panic("Unknown command: " + cmd)
	}
}
