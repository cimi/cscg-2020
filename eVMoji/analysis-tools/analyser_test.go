package main

import (
	"bytes"
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

func TestProgramDecode(t *testing.T) {
	base := "resources/"
	p := load(base + "code.bin")
	p.write(base + "reserialised.bin")

	start, _ := ioutil.ReadFile(base + "code.bin")
	reser, _ := ioutil.ReadFile(base + "reserialised.bin")
	if !bytes.Equal(start, reser) {
		t.Errorf("Reserialisation does not match original")
	}

	ioutil.WriteFile(base+"code.decoded.txt", []byte(decodeInstrs(p)), 0644)
	p.code = encodeFile(base + "code.decoded.txt")
	p.write(base + "reencoded.bin")
	reenc, _ := ioutil.ReadFile(base + "reencoded.bin")

	q := load(base + "reencoded.bin")
	ioutil.WriteFile(base+"code.decoded.again.txt", []byte(decodeInstrs(q)), 0644)
	if !bytes.Equal(reenc, reser) {
		t.Errorf("Encoding does not match original")
	}
}

func TestInstructionEncode(t *testing.T) {
	encoded := "I001824WD*131972!+*001624>101010+T142436>001010*001235S"
	decoded := "I25WD*140!+*23>0+T236>1*128S"
	i1 := decodeInstr(encoded)
	if i1 != decoded {
		t.Errorf("Decoded instruction does not match original: %s %s", i1, decoded)
	}

	i2 := encodeInstr(decoded)
	if string(i2.plain) != encoded {
		t.Errorf("Re-encoded instruction does not match original: %s %s", i2.plain, encoded)
	}
}
