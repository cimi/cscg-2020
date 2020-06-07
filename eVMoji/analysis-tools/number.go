package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

const encodingFile = "resources/number-encoding.json"
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
