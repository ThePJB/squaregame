package main

import (
	"flag"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Context struct {
	tStart  time.Time
	logfile *os.File
}

var gc GraphicsContext
var c Context

func main() {
	{
		// Setup
		c.tStart = time.Now()
		c.logfile = os.Stdout
		rand.Seed(time.Now().UnixNano())
		runtime.LockOSThread()

		// Flags
		xres := flag.Int("xres", 900, "x resolution of game window")
		yres := flag.Int("yres", 900, "y resolution of game window")
		flag.Parse()

		// Game settings

		// Graphics settings
		gc.xres = int32(*xres)
		gc.yres = int32(*yres)
		gc.frameInterval = (1_000_000_000 / 60) * time.Nanosecond
		gc.atlas = map[string]*sdl.Texture{}
		log(1, "resolution:", gc.xres, gc.yres)

		// Init engine
		initSDL()
		defer teardownSDL()

		//initTileTypes()
		//initEntityTypes()

		running := true
		//tEnd := time.Now()

		for running == true {
			FrameTimeStart := time.Now()
			//dt := (FrameTimeStart.Sub(tEnd)).Seconds()

			// handle input
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch t := event.(type) {
				case *sdl.QuitEvent:
					log(1, "User quit", t)
					running = false
				}
			}

			gc.renderer.Clear()
			gc.renderer.SetDrawColor(0, 0, 0, 255)
			gc.renderer.FillRect(&sdl.Rect{0, 0, gc.xres, gc.yres})

			gc.Draw()
			gc.renderer.Present()

			//tEnd = time.Now()
			tss := time.Since(FrameTimeStart)
			dtss := gc.frameInterval - tss
			if dtss > 0 {
				time.Sleep(dtss)
			}
		}
	}

}
