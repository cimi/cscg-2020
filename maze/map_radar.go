package main

import (
	"encoding/binary"
	"regexp"
	"strconv"
	"strings"
)

func teleportHangarBasement() []byte {
	decoded := make([]byte, 14)
	decoded[0] = 0x54
	decoded[1] = 0x01
	binary.LittleEndian.PutUint32(decoded[2:6], 1030622)
	copy(decoded[6:10], []byte{0x48, 0x77, 0xff, 0xff})
	binary.LittleEndian.PutUint32(decoded[10:14], 333486)
	return decoded
}

var re = regexp.MustCompile(`.*=\s+([^(]+)`)

func getNum(s string) int {
	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return 0
	}
	num, err := strconv.Atoi(matches[0][1])
	if err != nil {
		panic(err)
	}
	return num
}

func getNums(s string) []int {
	parts := strings.Split(s, "|")
	result := []int{}
	for _, p := range parts {
		result = append(result, getNum(p))
	}
	return result
}

func toF(num int) float64 {
	return float64(num) / 10000.0
}
