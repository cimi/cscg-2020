package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func decodeInstrs(p *Program) string {
	// asert write works correctly
	txt := ""
	for _, instr := range p.code {
		txt += instr.decoded
		txt += "\n"
	}
	return txt
}

func init() {
	fmt.Println("Hello from init!")
}

func TestProgramDecode(t *testing.T) {
	// convert the original to decoded text form
	base := "analysis/"
	p := load(base + "code.bin")
	p.writeBinary(base + "reserialised.bin")

	start, _ := ioutil.ReadFile(base + "code.bin")
	reser, _ := ioutil.ReadFile(base + "reserialised.bin")
	if !bytes.Equal(start, reser) {
		t.Errorf("Reserialisation does not match original")
	}

	ioutil.WriteFile(base+"code.decoded.txt", []byte(decodeInstrs(p)), 0644)
	p.code = encodeFile(base + "code.decoded.txt")
	p.writeBinary(base + "reencoded.bin")
	reenc, _ := ioutil.ReadFile(base + "reencoded.bin")

	q := load(base + "reencoded.bin")
	ioutil.WriteFile(base+"code.decoded.again.txt", []byte(decodeInstrs(q)), 0644)
	if !bytes.Equal(reenc, reser) {
		t.Errorf("Encoding does not match original")
	}
	// assert decode/encode works correctly

	// then convert back and check bytes
}

func TestInstructionEncode(t *testing.T) {
	encoded := "I001824WD*131972!+*001624>101010+T142436>001010*001235S"
	decoded := "I25WD*140!+*23>0+T236>1*128S"
	i1 := instr(encoded)
	if i1.decoded != decoded {
		t.Errorf("Decoded instruction does not match original: %s %s", i1.decoded, decoded)
	}

	i2 := instr(c.encodeInstr(decoded))
	if string(i2.plain) != encoded {
		t.Errorf("Re-encoded instruction does not match original: %s %s", i2.plain, encoded)
	}
}
