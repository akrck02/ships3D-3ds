package pixelart

import (
	"testing"

	"github.com/akrck02/go-renderer/draw2d"
)

func TestScalingMakesWhatWasDrawnLargerFromTheTopLeft(test *testing.T) {
	var drawing draw2d.List
	drawing.Rectangle(draw2d.Rect{X: 1, Y: 2, Width: 3, Height: 4}, 0, draw2d.Hex("#ffffff"), draw2d.Color{}, 0)
	drawing.Text("hi", 5, 6, 8, draw2d.Mono, draw2d.Hex("#ffffff"))
	Scale(&drawing, 1, 2)
	if got := drawing.Commands[0].Area; got != (draw2d.Rect{X: 1, Y: 2, Width: 3, Height: 4}) {
		test.Fatalf("the command before start moved: %v", got)
	}
	text := drawing.Commands[1]
	if text.Area.X != 10 || text.Baseline != 12 || text.Size != 16 {
		test.Fatalf("the text is at %v, baseline %v, size %v", text.Area.X, text.Baseline, text.Size)
	}
}
