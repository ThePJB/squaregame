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
// Levels get cluttered very quickly. maybe you want some dead squares that dont contribute to clutter
// Maybe cleaner graphics, circle tiles? or the colour of behav symbol/tile itself shows the state
// start with swaps and 1 tile things

// Puzzle elements:
// re cleaning: and gate, with the 4 tile toggle thing

// todos:
// level editor
// 	base config, reset
// clean up appearance
// 	always blank (doesnt matter) tiles
// 	tile state shines through symbol

// support resize for editor clarity

type Context struct {
	tStart  time.Time
	logfile *os.File

	level      level
	behavTypes []behavType

	mode            int
	editUISelection int
}

const (
	MODE_PLAY = iota
	MODE_EDIT
	NUM_MODES
)

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
		c.level.SetCellBehav(1, 1, BEHAV_SWAP)
		c.level.SetCellBehav(2, 2, BEHAV_SWAP)
		c.level.SetCellBehav(3, 3, BEHAV_SWAP)
		c.level.SetCellBehav(0, 1, BEHAV_TOGGLE)
		c.level.SetCellBehav(1, 0, BEHAV_TOGGLE)
		c.level.SetCellBehav(3, 3, BEHAV_SWAP)
		c.level.SetCellBehav(3, 4, BEHAV_SWAP)
		c.level.SetCellBehav(4, 4, BEHAV_SWAP)
		c.level.SetCellBehav(5, 5, BEHAV_SWAP)
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
				case *sdl.KeyboardEvent:
					if t.State == sdl.PRESSED {
						if t.Keysym.Sym == sdl.K_e {
							c.mode = (c.mode + 1) % NUM_MODES
						} else {
							for i := range c.behavTypes {
								if t.Keysym.Sym == c.behavTypes[i].edHotkey {
									c.editUISelection = i
									log(1, "Editor select", i, c.behavTypes[i].name)
									break
								}
							}
						}
					}
				}
			}

			gc.renderer.Clear()
			gc.renderer.SetDrawColor(64, 0, 64, 255)
			gc.renderer.FillRect(&sdl.Rect{0, 0, gc.xres, gc.yres})

			//gc.Draw()

			c.level.Draw(gc.xres, gc.yres)

			if c.mode == MODE_EDIT {
				drawEditUI(gc.xres, gc.yres)
			}

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

func drawEditUI(w, h int32) {
	s := h / NUM_BEHAVS
	for i := 0; i < NUM_BEHAVS; i++ {
		spriteSize := s / 2
		spriteToRect := sdl.Rect{50, int32(i) * s, spriteSize, spriteSize}

		borderSize := (s * 6) / 10
		borderOffset := (borderSize - spriteSize) / 2

		borderToRect := sdl.Rect{50 - borderOffset, int32(i)*s - borderOffset, borderSize, borderSize}
		if c.editUISelection == i {
			gc.renderer.SetDrawColor(255, 0, 0, 255)
		} else {
			gc.renderer.SetDrawColor(255, 255, 255, 255)
		}

		gc.renderer.FillRect(&borderToRect)
		gc.renderer.CopyEx(c.behavTypes[i].texture, nil, &spriteToRect, 0.0, nil, sdl.FLIP_NONE)
	}
}
