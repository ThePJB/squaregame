package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

func log(level int, args ...interface{}) {
	pargs := make([]interface{}, 0)
	pargs = append(pargs, "[", level, "][", time.Now(), "][", time.Since(c.tStart), "]")
	pargs = append(pargs, args...)
	fmt.Fprintln(c.logfile, pargs...)
}

func initSDL() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	if imgflags := img.Init(img.INIT_PNG); imgflags != img.INIT_PNG {
		panic("failed to init png loading")
	}

	if err := mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		panic(err)
	}
	window, err := sdl.CreateWindow("Square Game", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		gc.xres, gc.yres, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	gc.window = window

	renderer, err := sdl.CreateRenderer(gc.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	gc.renderer = renderer
	gc.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
}

func teardownSDL() {
	fmt.Print("Tearing down SDL...")

	for _, v := range gc.atlas {
		v.Destroy()
	}

	gc.window.Destroy()
	gc.renderer.Destroy()
	img.Quit()
	mix.CloseAudio()
	mix.Quit()
	sdl.Quit()
	fmt.Println("Done")
}
