package main

import (
	emu "Go-Chip8/emulator"
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
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

	if sdlErr := sdl.Init(sdl.INIT_EVERYTHING); sdlErr != nil {
		panic(sdlErr)
	}

	defer sdl.Quit()

	window, winErr := sdl.CreateWindow("Chip 8 - "+filename, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		64*10, 32*10, sdl.WINDOW_SHOWN)

	if winErr != nil {
		panic(winErr)
	}

	defer window.Destroy()

	renderer, renErr := sdl.CreateRenderer(window, -1, 0)

	if renErr != nil {
		panic(renErr)
	}
	defer renderer.Destroy()

	for {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch et := event.(type) {
			case *sdl.QuitEvent:
				os.Exit(0)
			case *sdl.KeyboardEvent:
				if et.Type == sdl.KEYUP {
					switch et.Keysym.Sym {
					case sdl.K_1:
						c.SetKey(0x1, false)
					case sdl.K_2:
						c.SetKey(0x2, false)
					case sdl.K_3:
						c.SetKey(0x3, false)
					case sdl.K_4:
						c.SetKey(0xC, false)
					case sdl.K_q:
						c.SetKey(0x4, false)
					case sdl.K_w:
						c.SetKey(0x5, false)
					case sdl.K_e:
						c.SetKey(0x6, false)
					case sdl.K_r:
						c.SetKey(0xD, false)
					case sdl.K_a:
						c.SetKey(0x7, false)
					case sdl.K_s:
						c.SetKey(0x8, false)
					case sdl.K_d:
						c.SetKey(0x9, false)
					case sdl.K_f:
						c.SetKey(0xE, false)
					case sdl.K_z:
						c.SetKey(0xA, false)
					case sdl.K_x:
						c.SetKey(0x0, false)
					case sdl.K_c:
						c.SetKey(0xB, false)
					case sdl.K_v:
						c.SetKey(0xF, false)
					}
				} else if et.Type == sdl.KEYDOWN {
					switch et.Keysym.Sym {
					case sdl.K_1:
						c.SetKey(0x1, true)
					case sdl.K_2:
						c.SetKey(0x2, true)
					case sdl.K_3:
						c.SetKey(0x3, true)
					case sdl.K_4:
						c.SetKey(0xC, true)
					case sdl.K_q:
						c.SetKey(0x4, true)
					case sdl.K_w:
						c.SetKey(0x5, true)
					case sdl.K_e:
						c.SetKey(0x6, true)
					case sdl.K_r:
						c.SetKey(0xD, true)
					case sdl.K_a:
						c.SetKey(0x7, true)
					case sdl.K_s:
						c.SetKey(0x8, true)
					case sdl.K_d:
						c.SetKey(0x9, true)
					case sdl.K_f:
						c.SetKey(0xE, true)
					case sdl.K_z:
						c.SetKey(0xA, true)
					case sdl.K_x:
						c.SetKey(0x0, true)
					case sdl.K_c:
						c.SetKey(0xB, true)
					case sdl.K_v:
						c.SetKey(0xF, true)
					}
				}
			}
		}

		emu.Emulate(c) //opcode computation

		if c.DrawFlag {

			renderer.SetDrawColor(255, 0, 0, 255)
			renderer.Clear()

			for y := 0; y < 32; y++ {
				for x := 0; x < 64; x++ {
					if c.Gfx[y*64+x] == 1 {
						renderer.SetDrawColor(255, 255, 255, 255)
					} else {
						renderer.SetDrawColor(0, 0, 0, 255)
					}
					renderer.FillRect(&sdl.Rect{
						Y: int32(y * 10),
						X: int32(x * 10),
						W: 10,
						H: 10,
					})
				}
			}
			renderer.Present()
			c.DrawFlag = false

		}

		sdl.Delay(1000 / 60)
	}
}
