package main

import (
	emu "Go-Chip8/emulator"
	"fmt"
	"os"
	"time"
)

func main() {
	c := emu.NewChip8()

	var filename string

	if len(os.Args) < 2 {
		fmt.Println("Enter rom filename")
		return
	}

	filename = os.Args[1]

	err := emu.LoadROM(c, filename)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	timer := time.NewTicker(time.Second / 60)
	defer timer.Stop()

	for range timer.C {

		emu.Emulate(c)

		if c.DrawFlag {
			emu.DrawScreen(c)
			c.DrawFlag = false
		}

		if c.DelayTimer > 0 {
			c.DelayTimer--
		}

		if c.SoundTimer > 0 {
			c.SoundTimer--
		}
	}

}
