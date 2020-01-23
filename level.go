package main

import "github.com/veandco/go-sdl2/sdl"

type level struct {
	w, h  int32
	cells []cell
}

type cell struct {
	state     int
	behaviour int
}

const (
	STATE_BLACK = iota
	STATE_WHITE
	NUM_STATES
)

const (
	BEHAV_NONE = iota
	BEHAV_TOGGLE
	BEHAV_4TOGGLE
	BEHAV_SWAP
	BEHAV_4TRIGGER
	NUM_BEHAVS
)

func makeLevel(w, h int32) level {
	l := level{
		w:     w,
		h:     h,
		cells: []cell{},
	}
	for i := 0; i < int(w*h); i++ {
		l.cells = append(l.cells, cell{STATE_BLACK, BEHAV_NONE})
	}
	return l
}

func (l *level) SetCellBehav(x, y int32, b int) {
	l.cells[y*l.w+x].behaviour = b
}
func (l *level) GetCellBehav(x, y int32) int {
	return l.cells[y*l.w+x].behaviour
}
func (l *level) SetCellState(x, y int32, s int) {
	l.cells[y*l.w+x].state = s
}
func (l *level) GetCellState(x, y int32) (int, bool) {
	if x >= l.w || y >= l.h {
		return 0, false
	}
	return l.cells[y*l.w+x].state, true
}

// Click function
func (l *level) DoCellAction(x, y int32) {
	b := l.cells[y*l.w+x].behaviour
	switch b {
	case BEHAV_NONE:
		// pass
	case BEHAV_TOGGLE:
		l.ToggleCellState(x, y)
	case BEHAV_4TOGGLE:
		l.ToggleCellState(x+1, y)
		l.ToggleCellState(x-1, y)
		l.ToggleCellState(x, y+1)
		l.ToggleCellState(x, y-1)
	case BEHAV_SWAP:
		{
			tmp1, ok1 := l.GetCellState(x-1, y)
			tmp2, ok2 := l.GetCellState(x+1, y)
			if ok1 && ok2 {
				l.SetCellState(x+1, y, tmp1)
				l.SetCellState(x-1, y, tmp2)
			}
		}
		{
			tmp1, ok1 := l.GetCellState(x, y+1)
			tmp2, ok2 := l.GetCellState(x, y-1)
			if ok1 && ok2 {
				l.SetCellState(x, y-1, tmp1)
				l.SetCellState(x, y+1, tmp2)
			}
		}
	case BEHAV_4TRIGGER:
		if l.GetCellBehav(x+1, y) != BEHAV_4TRIGGER {
			l.DoCellAction(x+1, y)
		}
		if l.GetCellBehav(x-1, y) != BEHAV_4TRIGGER {
			l.DoCellAction(x-1, y)
		}
		if l.GetCellBehav(x, y-1) != BEHAV_4TRIGGER {
			l.DoCellAction(x, y-1)
		}
		if l.GetCellBehav(x, y+1) != BEHAV_4TRIGGER {
			l.DoCellAction(x, y+1)
		}
	}
}

// Does nothing if x,y is OOB
func (l *level) ToggleCellState(x, y int32) {
	if x >= l.w || y >= l.h {
		return
	}
	l.cells[y*l.w+x].state = (l.cells[y*l.w+x].state + 1) % NUM_STATES
}

func (l level) Draw(w, h int32) {
	cw := w / l.w
	ch := h / l.h
	halfSpace := int32(4)
	stateSpace := int32(32)

	for i := range l.cells {
		x := int32(i) % l.w
		y := int32(i) / l.w

		behavToRect := sdl.Rect{x*cw + halfSpace, y*ch + halfSpace, cw - 2*halfSpace, ch - 2*halfSpace}
		stateToRect := sdl.Rect{x*cw + stateSpace, y*ch + stateSpace, cw - 2*stateSpace, ch - 2*stateSpace}

		gc.renderer.CopyEx(c.behavTypes[l.GetCellBehav(x, y)].texture, nil, &behavToRect, 0.0, nil, sdl.FLIP_NONE)
		state, _ := l.GetCellState(x, y)
		if state == STATE_BLACK {
			gc.renderer.SetDrawColor(0, 0, 0, 255)
		} else {
			gc.renderer.SetDrawColor(255, 255, 255, 255)
		}
		gc.renderer.FillRect(&stateToRect)

	}

}

type behavType struct {
	name    string
	texture *sdl.Texture
	// action func ...
}

func initCells() {
	c.behavTypes = append(c.behavTypes, behavType{"none", loadTexture(texturePath("none"))})
	c.behavTypes = append(c.behavTypes, behavType{"toggle", loadTexture(texturePath("toggle"))})
	c.behavTypes = append(c.behavTypes, behavType{"4toggle", loadTexture(texturePath("4toggle"))})
	c.behavTypes = append(c.behavTypes, behavType{"swap", loadTexture(texturePath("swap"))})
	c.behavTypes = append(c.behavTypes, behavType{"4trigger", loadTexture(texturePath("4trigger"))})

}
