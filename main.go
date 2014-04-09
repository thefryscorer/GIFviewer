package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
	"strconv"

	"github.com/banthar/Go-SDL/sdl"
)

func loadGIF(filepath string) *gif.GIF {
	infile, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Could not open %v : %v \n", filepath, err)
		os.Exit(1)
	}
	defer infile.Close()

	g, err := gif.DecodeAll(infile)
	if err != nil {
		fmt.Printf("Failed to decode GIF %v: %v \n", filepath, err)
		os.Exit(1)
	}

	return g
}

var (
	sdlScreen *sdl.Surface
)

func initDisplay(w, h int) error {
	sdlScreen = sdl.SetVideoMode(w, h, 32, sdl.SWSURFACE)

	if sdlScreen == nil {
		return errors.New("Error setting display mode.")
	}

	return nil
}

func Init(w, h int) error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err < 0 {
		return errors.New("sdl.Init failed with " + strconv.Itoa(err))
	}

	if err := initDisplay(w, h); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print("Too few arguments.")
		os.Exit(1)
	}
	g := loadGIF(os.Args[1])

	fmt.Printf("Loaded gif. Loopcount: %v. Number of frames: %v \n", g.LoopCount, len(g.Image))

	surfaces := make([]*sdl.Surface, 0)

	// Convert all images into sdl surfaces and add to list
	for n, i := range g.Image {
		fmt.Printf("Loading: %v/%v \n", n, len(g.Image))
		src := image.NewNRGBA(image.Rect(0, 0, i.Bounds().Max.X, i.Bounds().Max.Y))
		for x := 0; x < i.Bounds().Max.X; x++ {
			for y := 0; y < i.Bounds().Max.Y; y++ {
				src.Set(x, y, color.NRGBAModel.Convert(i.At(x, y)))
			}
		}
		surfaces = append(surfaces, sdl.CreateSurfaceFromImage(src))
	}

	defer sdl.Quit()

	// Initialise sdl display based on size of first image.
	w, h := g.Image[0].Bounds().Max.X, g.Image[0].Bounds().Max.Y
	err := Init(w, h)
	if err != nil {
		fmt.Printf("Could not init display: %v \n", err)
		os.Exit(1)
	}

	srcRect := sdl.Rect{
		W: uint16(w),
		H: uint16(h),
		X: 0,
		Y: 0,
	}
	if g.LoopCount == 0 {
		for {
			for i, s := range surfaces {
				sdlScreen.Blit(&srcRect, s, &srcRect)
				sdlScreen.Flip()
				sdl.Delay(uint32(g.Delay[i] * 10))
				e := sdl.PollEvent()
				if e != nil {
					if _, ok := e.(*sdl.QuitEvent); ok {
						sdl.Quit()
						os.Exit(0)
					}
				}
			}
		}

	} else {
		for i := 0; i > g.LoopCount; i++ {
			for i, s := range surfaces {
				sdlScreen.Blit(&srcRect, s, &srcRect)
				sdlScreen.Flip()
				sdl.Delay(uint32(g.Delay[i] * 10))
				e := sdl.PollEvent()
				if e != nil {
					if _, ok := e.(*sdl.QuitEvent); ok {
						sdl.Quit()
						os.Exit(0)
					}
				}
			}

		}
	}
}
