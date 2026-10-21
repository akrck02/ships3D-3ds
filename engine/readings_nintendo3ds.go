//go:build nintendo3ds

package main

import (
	"github.com/akrck02/go-engine/engine"
	"github.com/akrck02/go-renderer/ctr"
	"github.com/akrck02/go-renderer/draw2d"

	"github.com/jonbul/ships3D-3ds/game"
	"github.com/jonbul/ships3D-3ds/hud"
)

// showReadings paints the HUD over the top screen's 3D: the scene's own drawing, which Go UI's
// 3DS window paints in both eyes.
func showReadings(scene *engine.Scene, drawing *draw2d.List, state *game.Game, now float64) {
	hud.Paint(drawing, state, now, ctr.ScreenWidth, ctr.ScreenHeight)
	scene.World().Overlay.Drawing = drawing
}
