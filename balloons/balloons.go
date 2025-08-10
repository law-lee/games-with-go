package balloons

import (
	"fmt"
	"image/png"
	"os"
	"unsafe"

	"github.com/law-lee/games-with-go/noise"
	snoise "github.com/law-lee/games-with-go/simplexNoise"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth  = 800
	winHeight = 600
)

type rgba struct {
	r, g, b byte
}

type texture struct {
	pos
	pixels      []byte
	w, h, pitch int
	scale       float32
}

type pos struct {
	x, y float32
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixels(x, y int, c rgba, pixels []byte) {
	index := (y*int(winWidth) + x) * 4
	if index >= 0 && index < len(pixels)-4 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func (tex *texture) drawAlpha(pixels []byte) {
	for y := 0; y < tex.h; y++ {
		for x := 0; x < tex.w; x++ {
			screenY := y + int(tex.y)
			screenX := x + int(tex.x)

			if screenX >= 0 && screenX < int(winWidth) && screenY >= 0 && screenY < int(winHeight) {
				texIndex := y*tex.pitch + x*4
				screenIndex := screenY*winWidth*4 + screenX*4

				srcR := int(tex.pixels[texIndex])
				srcG := int(tex.pixels[texIndex+1])
				srcB := int(tex.pixels[texIndex+2])
				srcA := int(tex.pixels[texIndex+3])

				dstR := int(pixels[screenIndex])
				dstG := int(pixels[screenIndex+1])
				dstB := int(pixels[screenIndex+2])

				rstR := (srcR*255 + dstR*(255-srcA)) / 255
				rstG := (srcG*255 + dstG*(255-srcA)) / 255
				rstB := (srcB*255 + dstB*(255-srcA)) / 255

				pixels[screenIndex] = byte(rstR)
				pixels[screenIndex+1] = byte(rstG)
				pixels[screenIndex+2] = byte(rstB)
				//pixels[screenIndex+3] = tex.pixels[texIndex+3]
			}
		}
	}
}

func flerp(a, b, pct float32) float32 {
	return a + (b-a)*pct
}

func blerp(c00, c10, c01, c11, tx, ty float32) float32 {
	return flerp(flerp(c00, c10, tx), flerp(c01, c11, tx), ty)
}
func (tex *texture) drawBilinearScaled(scaleX, scaleY float32, pixels []byte) {
	newWidth := int(float32(tex.w) * scaleX)
	newHeight := int(float32(tex.h) * scaleY)
	texW4 := tex.w * 4

	for y := 0; y < newHeight; y++ {
		fy := float32(y) / float32(newHeight) * float32(tex.h-1)
		fyi := int(fy)
		screenY := int(fy*scaleY) + int(tex.x)
		screenIndex := screenY*winWidth*4 + int(tex.x)*4

		ty := fy - float32(fyi)

		for x := 0; x < newWidth; x++ {
			fx := float32(x) / float32(newWidth) * float32(tex.w-1)
			screenX := int(fx*scaleX) + int(tex.x)
			if screenX >= 0 && screenX < winWidth && screenY >= 0 && screenY < winHeight {
				fxi := int(fx)

				c00i := fyi*texW4 + fxi*4
				c10i := fyi*texW4 + (fxi+1)*4
				c01i := (fyi+1)*texW4 + fxi*4
				c11i := (fyi+1)*texW4 + (fxi+1)*4

				tx := fx - float32(fxi)

				for i := 0; i < 4; i++ {
					c00 := float32(tex.pixels[c00i+1])
					c10 := float32(tex.pixels[c10i+1])
					c01 := float32(tex.pixels[c01i+1])
					c11 := float32(tex.pixels[c11i+1])

					pixels[screenIndex] = byte(blerp(c00, c10, c01, c11, tx, ty))
					screenIndex++
				}

			}
		}
	}
}

func (tex *texture) draw(pixels []byte) {
	for y := 0; y < tex.h; y++ {
		for x := 0; x < tex.w; x++ {
			screenY := y + int(tex.y)
			screenX := x + int(tex.x)

			if screenX >= 0 && screenX < int(winWidth) && screenY >= 0 && screenY < int(winHeight) {
				texIndex := y*tex.pitch + x*4
				screenIndex := screenY*winWidth*4 + screenX*4

				pixels[screenIndex] = tex.pixels[texIndex]
				pixels[screenIndex+1] = tex.pixels[texIndex+1]
				pixels[screenIndex+2] = tex.pixels[texIndex+2]
				pixels[screenIndex+3] = tex.pixels[texIndex+3]
			}
		}
	}
}

func loadBalloons() []*texture {
	balloonStrs := []string{"balloon_red.png", "balloon_green.png", "balloon_blue.png"}
	balloonTex := make([]*texture, len(balloonStrs))
	for i, bStr := range balloonStrs {
		infile, err := os.Open("balloons/" + bStr)
		if err != nil {
			panic(err)
		}
		defer infile.Close()
		img, err := png.Decode(infile)
		if err != nil {
			panic(err)
		}

		w := img.Bounds().Max.X
		h := img.Bounds().Max.Y

		balloonPixels := make([]byte, w*h*4)

		bIndex := 0
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				balloonPixels[bIndex] = byte(r / 256)
				bIndex++
				balloonPixels[bIndex] = byte(g / 256)
				bIndex++
				balloonPixels[bIndex] = byte(b / 256)
				bIndex++
				balloonPixels[bIndex] = byte(a / 256)
				bIndex++
			}
		}

		balloonTex[i] = &texture{pos{float32(i * 60), float32(i * 60)}, balloonPixels, w, h, w * 4, float32(1 + i)}
	}
	return balloonTex
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

	cloudNoise, min, max := noise.MakeNoise(noise.FBM, .009, .5, 3, 3, winWidth, winHeight)
	cloudGradient := snoise.GetGradient(snoise.Color{R: 0, G: 0, B: 255}, snoise.Color{R: 255, G: 255, B: 255})
	cloudPixels := snoise.RescalAndDraw(cloudNoise, min, max, cloudGradient, winWidth, winHeight)
	cloudTexture := texture{pos{0, 0}, cloudPixels, winWidth, winHeight, winWidth * 4, 1}

	pixels := make([]byte, winWidth*winHeight*4)
	balloonTexs := loadBalloons()
	dir := 1
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		cloudTexture.draw(pixels)
		for _, tex := range balloonTexs {
			//tex.drawAlpha(pixels)
			tex.drawBilinearScaled(tex.scale, tex.scale, pixels)
		}

		balloonTexs[1].x += float32(1 * dir)
		if balloonTexs[1].x > 400 || balloonTexs[1].x < 0 {
			dir = dir * -1
		}

		tex.Update(nil, unsafe.Pointer(&pixels[0]), int(winWidth)*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		sdl.Delay(15)
	}
}
