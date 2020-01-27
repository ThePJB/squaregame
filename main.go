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

// for imgui input stuff needs to be globally accessible

// next todo: visual cleanup
// row and column blocks
// editor export

type Context struct {
	tStart  time.Time
	logfile *os.File

	level      level
	behavTypes []behavType

	mode            int
	editorSelection int
}

const (
	MODE_PLAY = iota
	MODE_EDIT_BEHAVIOUR
	MODE_EDIT_INITIAL_STATE
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
						if c.mode == MODE_EDIT_BEHAVIOUR {
							c.level.SetCellBehav(tx, ty, c.editorSelection)
						} else if c.mode == MODE_EDIT_INITIAL_STATE {
							c.level.SetCellInitialState(tx, ty, c.editorSelection)
							c.level.SetCellState(tx, ty, c.editorSelection)
						} else if c.mode == MODE_PLAY {
							c.level.DoCellAction(tx, ty)
						}
					}
				case *sdl.KeyboardEvent:
					if t.State == sdl.PRESSED {
						if t.Keysym.Sym == sdl.K_e {
							c.mode = (c.mode + 1) % NUM_MODES
							c.editorSelection = 0
						} else if c.mode != MODE_PLAY && t.Keysym.Sym >= sdl.K_1 && t.Keysym.Sym <= sdl.K_9 {
							c.editorSelection = int(t.Keysym.Sym - sdl.K_1)
						} else if t.Keysym.Sym == sdl.K_r {
							for i := range c.level.cells {
								c.level.cells[i].state = c.level.cells[i].initialState
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

			if c.mode == MODE_EDIT_BEHAVIOUR {
				drawSelectUI(60, 10, 20, NUM_BEHAVS, func(i int, toRect *sdl.Rect) {
					gc.renderer.CopyEx(c.behavTypes[i].texture, nil, toRect, 0, nil, sdl.FLIP_NONE)
				}, c.editorSelection)
			} else if c.mode == MODE_EDIT_INITIAL_STATE {
				drawSelectUI(60, 10, 20, NUM_STATES, func(i int, toRect *sdl.Rect) {
					if i == 0 {
						gc.renderer.SetDrawColor(255, 255, 255, 255)
					} else {
						gc.renderer.SetDrawColor(0, 0, 0, 255)
					}
					gc.renderer.FillRect(toRect)
				}, c.editorSelection)
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

// factor into draw numbered menu: individual s, border size, gap, images, callbacks
// actually just return selection

// actually idk if this immediate mode is quite how sdl is set up to work currently
// I wonder if there is a guarantee with SDL that inputs will persist for a frame

func drawSelectUI(sideLength, borderSize, gap int32, n int, draw func(i int, toRect *sdl.Rect), currentSelection int) {
	for i := 0; i < n; i++ {
		offset := gap + (gap+sideLength+2*borderSize)*int32(i)
		borderToRect := &sdl.Rect{gap, offset, sideLength + 2*borderSize, sideLength + 2*borderSize}
		imageToRect := &sdl.Rect{gap + borderSize, offset + borderSize, sideLength, sideLength}
		if i == currentSelection {
			gc.renderer.SetDrawColor(255, 0, 0, 255)
		} else {
			gc.renderer.SetDrawColor(128, 128, 128, 255)
		}
		gc.renderer.FillRect(borderToRect)
		draw(i, imageToRect)
	}
}
