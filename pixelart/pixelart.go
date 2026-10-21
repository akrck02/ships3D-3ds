// Package pixelart draws the game's pixel art in a draw2d list, as the 3DS's menus and HUD draw
// it: rectangles with their corners cut in steps, pictures drawn by hand in text, things that
// blink. Sizes are in pixels.
package pixelart

import (
	"math"

	"github.com/akrck02/go-renderer/draw2d"
)

// CornerStep is each step of a corner's staircase.
const CornerStep = 2

// TextSize is the pixel font's (VT323) size for text, and LabelSize the size for labels
// (Silkscreen and Press Start 2P), which are crisp at multiples of 8.
const (
	TextSize  = 16
	LabelSize = 8
)

// Stepped fills area with its corners cut in a staircase of steps steps: one rectangle a step,
// the widest in the middle.
func Stepped(drawing *draw2d.List, area draw2d.Rect, steps int, colour draw2d.Color) {
	for step := range steps + 1 {
		inset := float32((steps - step) * CornerStep)
		height := float32(step * CornerStep)
		drawing.Rectangle(draw2d.Rect{X: area.X + height, Y: area.Y + inset, Width: area.Width - 2*height, Height: area.Height - 2*inset}, 0, colour, draw2d.Color{}, 0)
	}
}

// Inset is area shrunk by a margin on every side.
func Inset(area draw2d.Rect, margin float32) draw2d.Rect {
	return draw2d.Rect{X: area.X + margin, Y: area.Y + margin, Width: area.Width - 2*margin, Height: area.Height - 2*margin}
}

// Shown tells whether something that blinks blinksASecond times a second shows at now: it shows
// for the share onShare of each blink.
func Shown(now, blinksASecond, onShare float64) bool {
	return math.Mod(now*blinksASecond, 1) < onShare
}

// PaintedBitmap draws a picture drawn by hand in text, each letter a colour of palette and each
// pixel size pixels across; letters the palette lacks are left clear.
func PaintedBitmap(drawing *draw2d.List, rows []string, x, y, size float32, palette map[byte]draw2d.Color) {
	for row, line := range rows {
		for column := 0; column < len(line); column++ {
			colour, found := palette[line[column]]
			if !found {
				continue
			}
			end := column + 1
			for end < len(line) && line[end] == line[column] {
				end++
			}
			drawing.Rectangle(draw2d.Rect{X: x + float32(column)*size, Y: y + float32(row)*size, Width: float32(end-column) * size, Height: size}, 0, colour, draw2d.Color{}, 0)
			column = end - 1
		}
	}
}
