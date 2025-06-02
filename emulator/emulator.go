package emulator

import "fmt"

func Emulate(c *Chip8) {
	opcode := (uint16(c.memory[c.pc]) << 8) | uint16(c.memory[c.pc+1])

	fmt.Printf("Opcode: 0x%04X at PC: 0x%03X\n", opcode, c.pc)

	switch opcode & 0xF000 {

	case 0x0000:
		switch opcode & 0x000F {
		case 0x0000:
			c.Gfx = [64 * 32]byte{}

		case 0x000E:
			c.sp--
			c.pc = c.stack[c.sp]
			return

		}

	case 0x1000:
		c.pc = opcode & 0x0FFF
		return

	case 0x2000:
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = opcode & 0x0FFF
		return

	case 0x3000:
		x := (opcode & 0x0F00) >> 8
		if c.V[x] == byte(opcode&0x00FF) {
			c.pc += 2
		}

	case 0x4000:
		x := (opcode & 0x0F00) >> 8
		if c.V[x] != byte(opcode&0x00FF) {
			c.pc += 2
		}

	case 0x5000:
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		if c.V[x] == c.V[y] {
			c.pc += 2
		}

	case 0x6000:
		x := (opcode & 0x0F00) >> 8
		c.V[x] = byte(opcode & 0x00FF)

	case 0x7000:
		x := (opcode & 0x0F00) >> 8
		c.V[x] += byte(opcode & 0x00FF)

	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			c.V[x] = c.V[y]
		case 0x0001:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			c.V[x] |= c.V[y]
		case 0x0002:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			c.V[x] &= c.V[y]
		case 0x0003:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			c.V[x] ^= c.V[y]
		case 0x0004:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			sum := uint16(c.V[x]) + uint16(c.V[y])
			c.V[0xF] = 0x0
			if sum > 255 {
				c.V[0xF] = 0x1
			}
			c.V[x] = byte(sum)
		case 0x0005:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			c.V[0xF] = 0x1
			if c.V[x] < c.V[y] {
				c.V[0xF] = 0x0
			}
			c.V[x] -= c.V[y]
		case 0x0006:
			x := (opcode & 0x0F00) >> 8
			bit := c.V[x] & 0x01
			c.V[0xF] = bit
			c.V[x] >>= 1
		case 0x0007:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			c.V[0xF] = 0x1
			if c.V[y] < c.V[x] {
				c.V[0xF] = 0x0
			}
			c.V[x] = c.V[y] - c.V[x]
		case 0x000E:
			x := (opcode & 0x0F00) >> 8
			bit := (c.V[x] & 0x80) >> 7
			c.V[0xF] = bit
			c.V[x] <<= 1
		}
	case 0x9000:
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		if c.V[x] != c.V[y] {
			c.pc += 2
		}

	case 0xA000:
		c.I = opcode & 0x0FFF

	case 0xB000:
		c.pc = (opcode & 0x0FFF) + uint16(c.V[0])
		return

	case 0xC000:
		x := (opcode & 0x0F00) >> 8
		rand := GetRand()
		NN := byte(opcode & 0x00FF)
		c.V[x] = rand & NN

	case 0xD000:
		vx := c.V[(opcode&0x0F00)>>8]
		vy := c.V[(opcode&0x00F0)>>4]
		n := int(opcode & 0x000F)

		c.V[0xF] = 0x0
		for i := 0; i < n; i++ {
			spriteByte := c.memory[c.I+uint16(i)] //example byte: 00111100
			for j := 0; j < 8; j++ {
				pixel := spriteByte >> (7 - j) & 1 //get a single bit from the byte
				if pixel == 0x01 {
					X := (int(vx) + j) % 64
					Y := (int(vy) + i) % 32
					index := Y*64 + X
					if c.Gfx[index] == 1 {
						c.V[0xF] = 0x1
					}
					c.Gfx[index] ^= 1
				}
			}
		}
		c.DrawFlag = true

	case 0xE000:
		switch opcode & 0x000F {
		case 0x000E:
			x := (opcode & 0x0F00) >> 8
			if c.key[c.V[x]] == 1 {
				c.pc += 2
			}
		case 0x0001:
			x := (opcode & 0x0F00) >> 8
			if c.key[c.V[x]] == 0 {
				c.pc += 2
			}
		}

	case 0xF000:
		switch opcode & 0x000F {
		case 0x0007:
			x := (opcode & 0x0F00) >> 8
			c.V[x] = byte(c.DelayTimer)
		case 0x000A:
			x := (opcode & 0x0F00) >> 8
			for i := 0; i < 16; i++ {
				if c.key[i] == 1 {
					c.V[x] = c.key[i]
					return
				}
			}
			c.pc -= 2
		case 0x0005:
			x := (opcode & 0x0F00) >> 8
			c.DelayTimer = int(c.V[x])
		case 0x0008:
			x := (opcode & 0x0F00) >> 8
			c.SoundTimer = int(c.V[x])

		case 0x000E:
			x := (opcode & 0x0F00) >> 8
			c.I += uint16(c.V[x])

		case 0x0009:
			x := (opcode & 0x0F00) >> 8
			c.I = 0x50 + uint16(c.V[x])*5
		case 0x0003:
			x := (opcode & 0x0F00) >> 8
			vx := int(c.V[x])
			c.memory[c.I] = byte(vx / 100)
			c.memory[c.I+1] = byte(vx / 10 % 10)
			c.memory[c.I+2] = byte(vx % 10)
		}
		switch opcode & 0x00F0 {
		case 0x0050:
			x := int((opcode & 0x0F00) >> 8)
			for i := 0; i <= x; i++ {
				c.memory[int(c.I)+i] = c.V[i]
			}
		case 0x0060:
			x := int((opcode & 0x0F00) >> 8)
			for i := 0; i <= x; i++ {
				c.V[i] = c.memory[int(c.I)+i]
			}
		}

	}

	c.pc += 2
}
