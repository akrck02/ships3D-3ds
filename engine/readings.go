//go:build !nintendo3ds

package main

import (
	"github.com/akrck02/go-engine/engine"
	"github.com/akrck02/go-renderer/draw2d"

	"github.com/jonbul/ships3D-3ds/game"
)

// showReadings would paint the HUD over the window's 3D view.
//
// TODO(akrck02): on the desktop Go UI's window paints its own interface as the scene's drawing,
// over anything put there; go-engine's 2D layer over 3D (its proposal 9, v0.0.5) is what lets
// the HUD show here, through hud.Paint as on the 3DS.
func showReadings(*engine.Scene, *draw2d.List, *game.Game, float64) {}
