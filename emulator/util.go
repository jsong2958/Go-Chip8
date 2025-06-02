package emulator

import (
	"math/rand"
	"os"
)

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
