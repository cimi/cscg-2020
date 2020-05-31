package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Program kills annoying linters.
type Program struct {
	data []byte
	code []*Instruction
}

func (p *Program) writeBinary(path string) {
	f, err := os.Create(path)
	defer f.Close()
	check(err)
	f.Write(p.data)
	for _, i := range p.code {
		f.Write(i.bytes)
	}
}

func (p *Program) writeDecoded(path string) {
	f, err := os.Create(path)
	defer f.Close()
	check(err)
	splitter := regexp.MustCompile(`[^\d]\d*`)
	for _, i := range p.code {
		f.WriteString(strings.Join(splitter.FindAllString(i.decoded, -1), "\n") + "\n")
	}
}

func (p *Program) pretty(text bool) string {
	offset := 0
	res := ""
	res += fmt.Sprintf("Data section:\n\n%3d -> %s\n%3d -> %s\n%3d -> %s\n%3d -> %s\n%3d -> %s\n",
		144, "Welcome to eVMoji 😎",
		167, "🤝 me the 🏳️",
		187, "tRy hArder! 💀💀💀",
		212, "Gotta go cyclic ♻️",
		235, "Thats the flag: CSCG{}", // 257
	)
	res += fmt.Sprintf("Code section:\n\n")
	for idx, i := range p.code {
		if text {
			res += fmt.Sprintf("%3d | %5d | %5d | %30s | %s\n", idx, offset, offset+512, i.decoded, i.plain)
		} else {
			res += fmt.Sprintf("%3d | %5d | %5d | %s\n", idx, offset, offset+512, i.bytes)
		}
		offset += len(i.bytes)
	}
	return res
}

func (p *Program) execute() string {
	tmpFile := "tmp.bin"
	p.writeBinary(tmpFile)
	cmd := exec.Command("./eVMoji", tmpFile)
	out, _ := cmd.CombinedOutput()
	// return fmt.Sprintf("%d", len(out))
	return string(out)
}

func load(binfilePath string) *Program {
	p := &Program{}
	binfile, err := ioutil.ReadFile(binfilePath)
	check(err)

	offset := 0x200
	p.data = binfile[0:offset]
	p.code = []*Instruction{}

	code := binfile[offset:]
	bs := []byte{}
	for utf8.RuneCount(code) >= 1 {
		r, size := utf8.DecodeRune(code)
		if r == muscle && len(bs) > 0 {
			p.code = append(p.code, NewInstruction(bs))
			bs = []byte{}
		}
		bs = append(bs, code[0:size]...)
		code = code[size:]
	}
	p.code = append(p.code, NewInstruction(bs))
	return p
}

/*
	// the number of bytes written leaks how the numbers are interpreted
	columns := []int{0, 1, 2, 3, 4, 10, 11, 12, 13, 14, 20, 21, 22, 23, 24, 31, 32, 33, 34}
	fmt.Printf("00xxxx:\t")
	for _, j := range columns {
		fmt.Printf("%7s,", fmt.Sprintf("00%02dxx", j))
	}
	fmt.Println()

	// we figure out that addresses follow the formula in decode()
	for i := 0; i <= 64; i++ {
		fmt.Printf("33xx%02d:\t", i)
		for _, j := range columns {
			num := fmt.Sprintf("33%02d%02d", j, i)
			p.code = []*Instruction{
				instr("I132672"),
				instr(fmt.Sprintf("I%sWD", num)),
			}
			l := p.execute()
			fmt.Printf("%7s,", l)
			ll, _ := strconv.Atoi(l)
			if decode(num) != ll {
				panic(fmt.Sprintf("%d != %d for %s", decode(num), ll, num))
			}
		}
		fmt.Println()
	}
*/
