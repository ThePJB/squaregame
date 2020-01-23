package main

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// purpose of this is to 'draw later'

type DrawRectCommand struct {
	x, y, w, h float64
	r, g, b, a float64
}

type GraphicsContext struct {
	window        *sdl.Window
	renderer      *sdl.Renderer
	atlas         map[string]*sdl.Texture
	xres, yres    int32
	frameInterval time.Duration

	DrawBuffer []DrawRectCommand
}

const (
	MODE_MARCH = iota
	MODE_HEIGHT
	MODE_BASIC
	NUM_GMODES
)

func (gc *GraphicsContext) DrawHArrow(x, y, m, w, r, g, b, a float64) {
	w = w / float64(gc.yres)
	if m > 0 {
		gc.DrawRect(x, y-w/2, m, w,
			r, g, b, a)
	} else {
		gc.DrawRect(x+m, y-w/2, -m, w,
			r, g, b, a)
	}
}

func (gc *GraphicsContext) DrawVArrow(x, y, m, w, r, g, b, a float64) {
	w = w / float64(gc.yres)
	if m > 0 {
		gc.DrawRect(x-w/2, y, w, m,
			r, g, b, a)
	} else {
		gc.DrawRect(x-w/2, y+m, w, -m,
			r, g, b, a)
	}
}

func (gc *GraphicsContext) DrawPoint(x, y, rad, r, g, b, a float64) {
	gc.DrawRect(x-rad/(2*float64(gc.xres)), y-rad/(2*float64(gc.yres)), rad/float64(gc.xres), rad/float64(gc.yres), r, g, b, a)
}

func (gc *GraphicsContext) DrawRect(x, y, w, h, r, g, b, a float64) {
	gc.DrawBuffer = append(gc.DrawBuffer, DrawRectCommand{x, y, w, h, r, g, b, a})
}

func (gc *GraphicsContext) Draw() {
	for _, r := range gc.DrawBuffer {
		dst := &sdl.Rect{int32(r.x * float64(gc.xres)),
			int32(r.y * float64(gc.yres)),
			int32(r.w * float64(gc.xres)),
			int32(r.h * float64(gc.yres))}

		gc.renderer.SetDrawColor(uint8(255*r.r), uint8(255*r.g), uint8(255*r.b), uint8(255*r.a))
		gc.renderer.FillRect(dst)
	}
	gc.DrawBuffer = []DrawRectCommand{}
}
