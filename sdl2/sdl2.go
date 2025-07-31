package sdl2

import (
	"fmt"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth  int32 = 800
	winHeight int32 = 600
)

type color struct {
	r, g, b byte
}

func setPixels(x, y int, c color, pixels []byte) {
	index := (y*int(winWidth) + x) * 4
	if index >= 0 && index < len(pixels)-4 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}
func Run() {
	window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, winWidth, winHeight)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	for y := 0; y < int(winHeight); y++ {
		for x := 0; x < int(winWidth); x++ {
			setPixels(x, y, color{byte(x % 255), byte(y % 255), byte(y % 255)}, pixels)
		}
	}
	tex.Update(nil, unsafe.Pointer(&pixels[0]), int(winWidth)*4)
	renderer.Copy(tex, nil, nil)
	renderer.Present()
	sdl.Delay(3000)
}
