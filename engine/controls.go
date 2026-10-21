//go:build !nintendo3ds

package main

import (
	"github.com/akrck02/go-engine/engine"

	"github.com/jonbul/ships3D-3ds/flight"
)

// controlsOf reads the keys as the flight model's controls, as the browser does: the arrows (or
// the 3DS's cross and circle pad) for the nose and the turn, A and D to roll, W and S for speed,
// Space (the 3DS's A) to fire.
func controlsOf(keys *engine.Input) (flight.Controls, bool) {
	turn, nose := keys.Axis(engine.KeyLeft, engine.KeyRight), keys.Axis(engine.KeyDown, engine.KeyUp)
	if stick := keys.State(); stick != nil {
		turn, nose = turn+stick.StickX, nose-stick.StickY
	}
	return flight.Controls{
		Pitch:    nose,
		Yaw:      turn,
		Roll:     keys.Axis(engine.KeyA, engine.KeyD),
		Throttle: keys.Axis(engine.KeyS, engine.KeyW),
	}, keys.Down(engine.KeySpace)
}

// pausePressed tells whether the player opened or closed the pause menu: Escape or P.
func pausePressed(keys *engine.Input) bool {
	return keys.Pressed(engine.KeyEscape) || keys.Pressed(engine.KeyP)
}
