package main

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth  int32 = 800
	winHeight int32 = 600
)

var nums = [][]byte{
	{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{
		1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
}

type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius float32
	xv     float32
	yv     float32
	color
}

type paddle struct {
	pos
	w     float32
	h     float32
	speed float32
	color
}

func drawNumber(pos pos, color color, pixelSize int, num int, pixels []byte) {
	startX := int(pos.x) - (pixelSize*3)/2
	startY := int(pos.y) - (pixelSize*5)/2

	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY+pixelSize; y++ {
				for x := startX; x < startX+pixelSize; x++ {
					setPixels(x, y, color, pixels)
				}
			}
		}
		startX += pixelSize
		if (i+1)%3 == 0 {
			startY += pixelSize
			startX -= pixelSize * 3
		}
	}
}

func (b *ball) draw(pixels []byte) {
	for y := -b.radius; y < b.radius; y++ {
		for x := -b.radius; x < b.radius; x++ {
			if x*x+y*y < b.radius*b.radius {
				setPixels(int(b.x+x), int(b.y+y), b.color, pixels)
			}
		}
	}
}

func getCenter() pos {
	return pos{float32(winWidth) / 2, float32(winHeight) / 2}
}
func (b *ball) update(leftPaddle, rightPaddle *paddle, elapsedTime float32) {
	b.x += b.xv * elapsedTime
	b.y += b.yv * elapsedTime

	if b.y < 0 || b.y > float32(winHeight) {
		b.yv = -b.yv
	}
	if b.x < 0 || b.x > float32(winWidth) {
		b.pos = getCenter()
	}

	if b.x < leftPaddle.x+leftPaddle.w/2 {
		if b.y > leftPaddle.y-leftPaddle.h/2 && b.y < leftPaddle.y+leftPaddle.h/2 {
			b.xv = -b.xv
		}
	}
	if b.x > rightPaddle.x-rightPaddle.w/2 {
		if b.y > rightPaddle.y-rightPaddle.h/2 && b.y < rightPaddle.y+rightPaddle.h/2 {
			b.xv = -b.xv
		}
	}

}

func (paddle *paddle) draw(pixels []byte) {
	startX := paddle.x - paddle.w/2
	startY := paddle.y - paddle.h/2

	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			setPixels(int(startX)+x, int(startY)+y, paddle.color, pixels)
		}
	}
}

func (p *paddle) update(keyState []uint8, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		p.y -= p.speed * elapsedTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		p.y += p.speed * elapsedTime
	}
}

func (p *paddle) aiUpdate(ball *ball, elapsedTime float32) {
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

	// for y := 0; y < winHeight); y++ {
	// 	for x := 0; x < winWidth); x++ {
	// 		setPixels(x, y, color{byte(x % 255), byte(y % 255), byte(y % 255)}, pixels)
	// 	}
	// }
	// tex.Update(nil, unsafe.Pointer(&pixels[0]), winWidth)*4)
	// renderer.Copy(tex, nil, nil)
	// renderer.Present()

	player1 := paddle{pos{100, 100}, 20, 100, 300, color{255, 255, 255}}
	player2 := paddle{pos{700, 100}, 20, 100, 300, color{255, 255, 255}}
	ball := ball{pos{300, 300}, 20, 400, 400, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	var frameStart time.Time
	var elapsedTime float32
	// OSX requires that you consume events for windows to open and work properly
	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		clear(pixels)
		drawNumber(getCenter(), color{255, 255, 255}, 3, 2, pixels)
		player1.draw(pixels)
		player1.update(keyState, elapsedTime)

		player2.draw(pixels)
		player2.aiUpdate(&ball, elapsedTime)

		ball.draw(pixels)
		ball.update(&player1, &player2, elapsedTime)

		tex.Update(nil, unsafe.Pointer(&pixels[0]), int(winWidth)*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())

		//sdl.Delay(16)
		if elapsedTime < .005 {
			sdl.Delay(5 - uint32(elapsedTime*1000))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
		//fmt.Println(elapsedTime)
	}
}
