package main

import (
	"fmt"
	"testing"
)

func TestBotLogin(t *testing.T) {
	winRace()
	t.Error("done!")
}

func TestDumpRacePath(t *testing.T) {
	for idx, cmd := range raceCommands {
		fmt.Println(idx, cmd.ts, cmd.position)
	}
	t.Error("done!")
}

func TestPositionAdjustment(t *testing.T) {
	src := position{2805794, 0, 2365136}
	dst := position{2772149, 0, 2333746}
	next := src

	for moves := 0; moves < 6; moves++ {
		next = next.moveTowards(dst)
	}
	if !next.equals(dst) {
		t.Errorf("Failed to reach destination %v: %v", dst, next)
	}

}
