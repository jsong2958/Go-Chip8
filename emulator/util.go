package emulator

import (
	"fmt"
	"math/rand"
	"os"
)

var keyMap = map[rune]int{
	'1': 0x1, '2': 0x2, '3': 0x3, '4': 0xC,
	'q': 0x4, 'w': 0x5, 'e': 0x6, 'r': 0xD,
	'a': 0x7, 's': 0x8, 'd': 0x9, 'f': 0xE,
	'z': 0xA, 'x': 0x0, 'c': 0xB, 'v': 0xF,
}

func GetRand() byte {
	return byte(rand.Intn(256))
}

func LoadROM(c *Chip8, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	copy(c.memory[0x200:], data)
	return nil

}

func DrawScreen(c *Chip8) {
	fmt.Print("\033[H\033[2J")
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if c.gfx[y*64+x] == 1 {
				fmt.Print("▐")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
