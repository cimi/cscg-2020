package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGetNum(t *testing.T) {
	if getNum("ts=    156240(0x50620200)|") != 156240 {
		t.Error("Failed number conversion")
	}
}

func TestMapRadarCoords(t *testing.T) {
	file, err := os.Open("notes/radar.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Index(line, "0x3713ffff") == -1 {
			continue
		}
		// fmt.Println(line)
		// use regex to extract position and time from string
		//ts=    156240(0x50620200)|?=         0(0x00000000)|x=    304340(0xd4a40400)|y=    -35000(0x4877ffff)|z=    250231(0x77d10300)|?=       860(0x5c030000)|?c=   1097732(0x04c01000)|?=         0(0x00000000)|0|200(0xc800)|0(0x0000) (uid=    -60617(0x3713ffff))
		nums := getNums(line)
		ts, x, y, z := toF(nums[0]), toF(nums[2]), toF(nums[3]), toF(nums[4])
		fmt.Printf("%f,%f,%f,%f\n", ts, x, y, z)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// t.Error("Done")
}
