//go:build nintendo3ds

package main

import (
	"github.com/akrck02/go-renderer/ctr"
	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/input"
	ui "github.com/akrck02/go-ui"
)

// topTitle paints the title's name over the top screen.
type topTitle struct {
	context  *ui.Context
	renderer *ctr.FlatRenderer
	state    *input.State
}

// newTopTitle prepares a go-ui context for the top screen. It has a painter of its own: one
// painter for both screens would overwrite the vertices of one before the GPU drew them.
func newTopTitle(bottomScreen *ctr.BottomScreen, theme ui.Theme) (*topTitle, error) {
	renderer, err := ctr.NewFlatRenderer(bottomScreen)
	if err != nil {
		return nil, err
	}
	context, err := ui.New()
	if err != nil {
		return nil, err
	}
	context.Theme = theme
	state := input.NewState()
	state.FramebufferWidth, state.FramebufferHeight = ctr.ScreenWidth, ctr.ScreenHeight
	state.WindowWidth, state.WindowHeight = ctr.ScreenWidth, ctr.ScreenHeight
	return &topTitle{context: context, renderer: renderer, state: state}, nil
}

// paint lays a heading and a caption out and paints them over the top screen, drawn just before
// in one picture.
func (title *topTitle) paint(screen *ctr.Screen, heading, caption string) {
	title.context.Begin(title.state)
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Fixed(ctr.ScreenWidth), Height: ui.Fixed(ctr.ScreenHeight), AlignX: ui.Center, Gap: 2,
		Padding: ui.Padding{Top: 28}}
	title.context.Box(layout, func() {
		title.context.Text(heading, ui.TextStyle{Size: titleSize, Weight: draw2d.Bold})
		title.context.Caption(caption)
	})
	title.context.End()
	title.renderer.PaintOver(screen, title.context.Drawing())
}

// paintOver paints a draw list of the program's own over the top screen, on both eyes' pictures
// when the frame is in 3D, so that it shows at the screen's depth.
func (title *topTitle) paintOver(screen *ctr.Screen, drawing *draw2d.List, stereo bool) {
	if !stereo {
		title.renderer.PaintOver(screen, drawing)
		return
	}
	for _, eye := range []ctr.Eye{ctr.LeftEye, ctr.RightEye} {
		screen.DrawOn(eye)
		title.renderer.PaintOver(screen, drawing)
	}
}
