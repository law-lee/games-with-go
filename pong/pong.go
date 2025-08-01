package pong

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/law-lee/games-with-go/noise"
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
	score int
	color
}

type gameState int

const (
	START gameState = iota
	PLAY
)

var state = START

func flerp(b1 byte, b2 byte, pct float32) byte {
	return byte(float32(b1) + pct*(float32(b2)-float32(b1)))
}

func colorLerp(c1, c2 color, pct float32) color {
	return color{flerp(c1.r, c2.r, pct), flerp(c1.g, c2.g, pct), flerp(c1.b, c2.b, pct)}
}

func getGradient(c1, c2 color) []color {
	result := make([]color, 256)

	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = colorLerp(c1, c2, pct)
	}
	return result
}

func getDualGradient(c1, c2, c3, c4 color) []color {
	result := make([]color, 256)

	for i := range result {
		pct := float32(i) / float32(255)
		if pct < 0.5 {
			result[i] = colorLerp(c1, c2, pct*float32(2))
		} else {
			result[i] = colorLerp(c3, c4, pct*float32(1.5)-float32(0.5))
		}
	}
	return result
}

func clamp(min, max, v int) int {
	if v < min {
		v = min
	} else if v > max {
		v = max
	}
	return v
}

func rescalAndDraw(noise []float32, min, max float32, gradient []color, w, h int32) []byte {
	result := make([]byte, w*h*4)
	scale := 255.0 / (max - min)
	offset := min * scale

	for i := range noise {
		noise[i] = noise[i]*scale - offset
		c := gradient[clamp(0, 255, int(noise[i]))]
		p := i * 4
		result[p] = c.r
		result[p+1] = c.g
		result[p+2] = c.b
	}
	return result
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
	if b.x < 0 {
		rightPaddle.score++
		b.pos = getCenter()
		state = START
	} else if b.x > float32(winWidth) {
		leftPaddle.score++
		b.pos = getCenter()
		state = START
	}

	if b.x-b.radius < leftPaddle.x+leftPaddle.w/2 {
		if b.y > leftPaddle.y-leftPaddle.h/2 && b.y < leftPaddle.y+leftPaddle.h/2 {
			b.xv = -b.xv
			b.x = leftPaddle.x + leftPaddle.w/2.0 + b.radius
		}
	}
	if b.x+b.radius > rightPaddle.x-rightPaddle.w/2 {
		if b.y > rightPaddle.y-rightPaddle.h/2 && b.y < rightPaddle.y+rightPaddle.h/2 {
			b.xv = -b.xv
			b.x = rightPaddle.x - rightPaddle.w/2 - b.radius
		}
	}

}

func lerp(a, b float32, pct float32) float32 {
	return a + pct*(b-a)
}

func (paddle *paddle) draw(pixels []byte) {
	startX := paddle.x - paddle.w/2
	startY := paddle.y - paddle.h/2

	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			setPixels(int(startX)+x, int(startY)+y, paddle.color, pixels)
		}
	}

	numX := lerp(paddle.x, getCenter().x, 0.2)
	drawNumber(pos{numX, 35}, paddle.color, 10, paddle.score, pixels)
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
func Run() {
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

	player1 := paddle{pos{100, 100}, 20, 100, 300, 0, color{255, 255, 255}}
	player2 := paddle{pos{700, 100}, 20, 100, 300, 0, color{255, 255, 255}}
	ball := ball{pos{300, 300}, 20, 400, 400, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	noise, min, max := noise.MakeNoise(noise.FBM, .01, 0.2, 2, 3, int(winWidth), int(winHeight))
	gradient := getGradient(color{255, 0, 0}, color{0, 0, 0})
	noisePixels := rescalAndDraw(noise, min, max, gradient, winWidth, winHeight)

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
		// for i := range noisePixels {
		// 	pixels[i] = noisePixels[i]
		// }
		copy(pixels, noisePixels)
		drawNumber(getCenter(), color{255, 255, 255}, 3, 2, pixels)
		switch state {
		case PLAY:
			player1.update(keyState, elapsedTime)
			player2.aiUpdate(&ball, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)
		case START:
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3 {
					player1.score = 0
					player2.score = 0
				}

				state = PLAY
			}
		}
		player1.draw(pixels)

		player2.draw(pixels)

		ball.draw(pixels)

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
