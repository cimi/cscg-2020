package main

import (
	"bytes"
)

// Instruction baby.
type Instruction struct {
	bytes   []byte
	runes   []rune
	plain   []byte
	decoded string
}

// NewInstruction parses an instruction from a byte slice.
func NewInstruction(original []byte) *Instruction {
	plain := c.toText(original)
	emoji := c.toEmoji(plain)
	if bytes.Compare(emoji, original) != 0 {
		panic("Wrong conversion!")
	}
	return &Instruction{
		bytes:   original,
		runes:   c.toRunes(original),
		plain:   plain,
		decoded: c.decodeInstr(string(plain)),
	}
}

func instr(original string) *Instruction {
	emoji := c.toEmoji([]byte(original))
	return NewInstruction(emoji)
}
