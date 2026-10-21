//go:build !nintendo3ds

package main

import "github.com/jonbul/ships3D-3ds/game"

// radarScreen is the radar, which on the 3DS is on the bottom screen.
//
// TODO(akrck02): the desktop has no bottom screen; the browser shows the 3D radar in a corner,
// which needs a second view in go-engine's window.
type radarScreen struct{}

func openRadar() *radarScreen { return &radarScreen{} }

func (*radarScreen) show(*game.Game, float64) {}

func (*radarScreen) close() {}
