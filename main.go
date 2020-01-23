package main

import (
	"flag"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// just realised that with triggering, need to be careful of order and stuff. only for non commutative ops

type Context struct {
	tStart  time.Time
	logfile *os.File

	level      level
	behavTypes []behavType
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
		c.level = makeLevel(6, 6)
		c.level.SetCellBehav(2, 2, BEHAV_TOGGLE)
		c.level.SetCellBehav(2, 3, BEHAV_TOGGLE)
		c.level.SetCellBehav(3, 3, BEHAV_SWAP)
		c.level.SetCellBehav(3, 4, BEHAV_SWAP)
		c.level.SetCellBehav(4, 4, BEHAV_SWAP)
		c.level.SetCellBehav(2, 4, BEHAV_TOGGLE)

		// Graphics settings
		gc.xres = int32(*xres)
		gc.yres = int32(*yres)
		gc.frameInterval = (1_000_000_000 / 60) * time.Nanosecond
		gc.atlas = map[string]*sdl.Texture{}
		log(1, "resolution:", gc.xres, gc.yres)

		// Init engine
		initSDL()
		defer teardownSDL()

		initCells()

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
				case *sdl.MouseButtonEvent:
					if t.Button == sdl.BUTTON_LEFT && t.State == sdl.PRESSED {
						tx := t.X * c.level.w / gc.xres
						ty := t.Y * c.level.h / gc.yres
						c.level.DoCellAction(tx, ty)
					}
				}
			}

			gc.renderer.Clear()
			gc.renderer.SetDrawColor(64, 0, 64, 255)
			gc.renderer.FillRect(&sdl.Rect{0, 0, gc.xres, gc.yres})

			//gc.Draw()

			c.level.Draw(gc.xres, gc.yres)

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
