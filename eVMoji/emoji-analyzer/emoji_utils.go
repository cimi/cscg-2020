package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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

const encodingFile = "analysis/number-encoding.json"
const muscle = 128170

var encodingTable = make(map[string]string)
var c = &converter{}

func init() {
	_, err := os.Stat(encodingFile)
	if os.IsNotExist(err) {
		// if we don't already have the encoding table, we decode all possible number representation
		for i := 0; i < 1000000; i++ {
			encoded := fmt.Sprintf("%06d", i)
			decoded := c.decodeNum(encoded)
			encodingTable[decoded] = encoded
		}
		b, err := json.MarshalIndent(encodingTable, "", "  ")
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile(encodingFile, b, 0644)
		fmt.Println("Generated encoding table in file")
	}
	check(err)
	bs, err := ioutil.ReadFile(encodingFile)
	check(err)
	err = json.Unmarshal(bs, &encodingTable)
	check(err)
}

type converter struct {
}

func (c *converter) toText(emoji []byte) []byte {
	bs := make([]byte, len(emoji))
	copy(bs, emoji)
	// convert all digit emojis to characters
	bs = bytes.ReplaceAll(bs, []byte{0xef, 0xb8, 0x8f, 0xe2, 0x83, 0xa3}, []byte{})
	for k, v := range emojiMap {
		bs = bytes.ReplaceAll(bs, []byte(k), []byte(v))
	}

	return bs
}

func (c *converter) toEmoji(text []byte) []byte {
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

func (c *converter) toRunes(bs []byte) []rune {
	rs := []rune{}
	for utf8.RuneCount(bs) >= 1 {
		r, size := utf8.DecodeRune(bs)
		rs = append(rs, r)
		bs = bs[size:]
	}
	return rs
}

func (c *converter) decodeInstr(txt string) string {
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

func (c *converter) encodeInstr(txt string) string {
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

func (c *converter) decodeNum(num string) string {
	_, err := strconv.Atoi(num)
	if err != nil || len(num) != 6 {
		return ""
	}
	x1, _ := strconv.Atoi(string(num[0]))
	x2, _ := strconv.Atoi(string(num[1]))
	x3, _ := strconv.Atoi(string(num[2]))
	x4, _ := strconv.Atoi(string(num[3]))
	x5, _ := strconv.Atoi(string(num[4]))
	x6, _ := strconv.Atoi(string(num[5]))
	res := strconv.Itoa(pow(x2, x1) + pow(x4, x3) + pow(x6, x5))
	encodingTable[res] = num
	return res
}

func (c *converter) encodeNum(str string) string {
	_, err := strconv.Atoi(str)
	if err == nil {
		val, ok := encodingTable[str]
		if !ok {
			panic("Encoding not found for " + str)
		}
		return val
	}
	return str
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
