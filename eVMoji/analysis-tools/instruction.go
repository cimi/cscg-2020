package main

import (
	"bytes"
	"unicode/utf8"
)

var emojiMap = map[string]string{
	"💪":  "I", // push to stack
	"🤔":  "T", // branch
	"➕":  "+", // bitwise and 1
	"🌠":  "*", // *(four bytes at address)
	"🔀":  "S", // xor
	"📖":  "R", // read
	"🦾":  "M", // index into array
	"✅":  "C", // or
	"💀":  "D", // exit
	"‼️": "!", // push
	"✏️": "W", // print
	"➡️": ">", // right shift by (param)
}

// Instruction baby.
type Instruction struct {
	bytes   []byte
	runes   []rune
	plain   []byte
	decoded string
}

// NewInstruction parses an instruction from a byte slice.
func NewInstruction(original []byte) *Instruction {
	plain := toText(original)
	emoji := toEmoji(plain)
	if bytes.Compare(emoji, original) != 0 {
		panic("Wrong conversion!")
	}
	return &Instruction{
		bytes:   original,
		runes:   toRunes(original),
		plain:   plain,
		decoded: decodeInstr(string(plain)),
	}
}

func encodeInstr(instr string) *Instruction {
	emoji := toEmoji([]byte(encodeNums(instr)))
	return NewInstruction(emoji)
}

func encodeNums(txt string) string {
	ss := re2.FindAllString(txt, -1)
	out := ""
	for _, p := range ss {
		// encode n to abcdef where b^a + d^c + f^e == n
		encoded := c.encodeNum(p)
		if encoded != "" {
			out += encoded
		} else {
			out += p
		}
	}
	return out
}

func decodeInstr(txt string) string {
	ss := re2.FindAllString(txt, -1)
	out := ""
	for _, s := range ss {
		decoded := c.decodeNum(s)
		if decoded != "" {
			out += decoded
		} else {
			out += s
		}
	}
	return out
}

func toEmoji(text []byte) []byte {
	bs := make([]byte, len(text))
	copy(bs, text)
	// convert all digit characters to emoji
	for i := 0x30; i <= 0x39; i++ {
		bs = bytes.ReplaceAll(bs, []byte{byte(i)}, []byte{byte(i), 0xef, 0xb8, 0x8f, 0xe2, 0x83, 0xa3})
	}
	for k, v := range emojiMap {
		bs = bytes.ReplaceAll(bs, []byte(v), []byte(k))
	}
	return bs
}

func toText(emoji []byte) []byte {
	bs := make([]byte, len(emoji))
	copy(bs, emoji)
	// convert all digit emojis to characters
	bs = bytes.ReplaceAll(bs, []byte{0xef, 0xb8, 0x8f, 0xe2, 0x83, 0xa3}, []byte{})
	for k, v := range emojiMap {
		bs = bytes.ReplaceAll(bs, []byte(k), []byte(v))
	}

	return bs
}

func toRunes(bs []byte) []rune {
	rs := []rune{}
	for utf8.RuneCount(bs) >= 1 {
		r, size := utf8.DecodeRune(bs)
		rs = append(rs, r)
		bs = bs[size:]
	}
	return rs
}
