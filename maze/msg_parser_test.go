package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

type decodeTestCase struct {
	in  []byte
	out []byte
}

func TestEncodeDecode(t *testing.T) {
	expected := []byte{0x36, 0x78, 0x73, 0x33, 0x43, 0x8a, 0xb1, 0x3d, 0x9f, 0x93, 0x9f, 0x65}
	decoded := decode(expected)
	actual := encode(decoded, 0x36, 0x78)
	if bytes.Compare(expected, actual) != 0 {
		t.Errorf("Round trip failed: %x -> %x -> %x", expected, decoded, actual)
	}
}

func TestDecode(t *testing.T) {
	tcs := []*decodeTestCase{
		&decodeTestCase{
			in:  []byte{0x0, 0x0, 0x0, 0x0},
			out: []byte{0x0, 0x0},
		},
		&decodeTestCase{
			in:  []byte{0x36, 0x78, 0x73, 0x33, 0x43, 0x8a, 0xb1, 0x3d, 0x9f, 0x93, 0x9f, 0x65},
			out: []byte{0x45, 0x9D, 0x64, 0x15, 0xA9, 0xAD, 0x96, 0x12, 0x66, 0x17},
		},
		&decodeTestCase{
			in:  []byte{0x30, 0xcf, 0x75, 0x9d, 0xab, 0x8a, 0xc6, 0x92, 0x99, 0xcc, 0xc8, 0x69},
			out: []byte{0x45, 0x9D, 0x64, 0x15, 0xA9, 0xAD, 0x96, 0x12, 0x66, 0x17},
		},
	}
	for _, tc := range tcs {
		actual := decode(tc.in)
		expected := tc.out
		if bytes.Compare(actual, expected) != 0 {
			t.Errorf("Decode failed: %x -> %x != %x", tc.in, actual, expected)
		}
	}
}

type parseTestCase struct {
	in  []byte
	out *command
}

func TestParseClientMsg(t *testing.T) {
	tcs := []*parseTestCase{
		&parseTestCase{
			in: []byte{0x49, 0x9d, 0x64, 0x15, 0xa9, 0xad, 0x96, 0x12, 0x66, 0x37, 0x13, 0xff, 0xff},
			out: &command{
				code:   []byte{0x49},
				name:   "info",
				secret: secret,
				data:   []byte{0x37, 0x13, 0xff, 0xff},
				// parsed: "chars=[7  ÿ ÿ]",
			},
		},
	}
	for _, tc := range tcs {
		actual := parseClientMsg(tc.in)
		expected := tc.out
		fmt.Println(actual)
		deep.CompareUnexportedFields = true
		if diff := deep.Equal(expected, actual); diff != nil {
			t.Error(diff)
		}
	}
}

func TestParseClientCmdFile(t *testing.T) {
	file, err := os.Open("notes/decoded-msg.md")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Index(line, "[CLIENT len=") != 0 {
			continue
		}
		cmdBytes, err := hex.DecodeString(strings.Split(line, "] ")[1])
		if err != nil {
			panic(err)
		}
		cmd := parseClientMsg(cmdBytes)
		if cmd.name != "xxx" {
			fmt.Println(cmd)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	t.Error("Done")
}
