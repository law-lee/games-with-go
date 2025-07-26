package main

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

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius int
	xv     float32
	yv     float32
	color
}

type paddle struct {
	pos
	w int
	h int
	color
}

func (b *ball) draw(pixels []byte) {
	for y := -b.radius; y < b.radius; y++ {
		for x := -b.radius; x < b.radius; x++ {
			if x*x+y*y < b.radius*b.radius {
				setPixels(int(b.x)+x, int(b.y)+y, b.color, pixels)
			}
		}
	}
}

func (b *ball) update(leftPaddle, rightPaddle *paddle) {
	b.x += b.xv
	b.y += b.yv

	if b.y < 0 || int(b.y) > int(winHeight) {
		b.yv = -b.yv
	}
	if b.x < 0 || int(b.x) > int(winWidth) {
		b.x = 300
		b.y = 300
	}

	if int(b.x) < int(leftPaddle.x)+leftPaddle.w/2 {
		if int(b.y) > int(leftPaddle.y)-leftPaddle.h/2 && int(b.y) < int(leftPaddle.y)+leftPaddle.h/2 {
			b.xv = -b.xv
		}
	}
	if int(b.x) > int(rightPaddle.x)-rightPaddle.w/2 {
		if int(b.y) > int(rightPaddle.y)-rightPaddle.h/2 && int(b.y) < int(rightPaddle.y)+rightPaddle.h/2 {
			b.xv = -b.xv
		}
	}

}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x) - paddle.w/2
	startY := int(paddle.y) - paddle.h/2

	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixels(startX+x, startY+y, paddle.color, pixels)
		}
	}
}

func (p *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		p.y -= 5
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		p.y += 5
	}
}

func (p *paddle) aiUpdate(ball *ball) {
	p.y = ball.y
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}
func setPixels(x, y int, c color, pixels []byte) {
	index := (y*int(winWidth) + x) * 4
	if index >= 0 && index < len(pixels)-4 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}
func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
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

	// for y := 0; y < int(winHeight); y++ {
	// 	for x := 0; x < int(winWidth); x++ {
	// 		setPixels(x, y, color{byte(x % 255), byte(y % 255), byte(y % 255)}, pixels)
	// 	}
	// }
	// tex.Update(nil, unsafe.Pointer(&pixels[0]), int(winWidth)*4)
	// renderer.Copy(tex, nil, nil)
	// renderer.Present()

	player1 := paddle{pos{100, 100}, 20, 100, color{255, 255, 255}}
	player2 := paddle{pos{700, 100}, 20, 100, color{255, 255, 255}}
	ball := ball{pos{300, 300}, 20, 3, 3, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()
	// OSX requires that you consume events for windows to open and work properly
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		clear(pixels)
		player1.draw(pixels)
		player1.update(keyState)

		player2.draw(pixels)
		player2.aiUpdate(&ball)

		ball.draw(pixels)
		ball.update(&player1, &player2)

		tex.Update(nil, unsafe.Pointer(&pixels[0]), int(winWidth)*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		sdl.Delay(16)
	}
}
