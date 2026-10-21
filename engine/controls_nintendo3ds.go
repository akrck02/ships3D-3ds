//go:build nintendo3ds

package main

import (
	"github.com/akrck02/go-engine/engine"

	"github.com/jonbul/ships3D-3ds/flight"
)

// controlsOf reads the console as the 3DS program does: the circle pad for the nose (down raises
// it) and the turn, the cross's left and right (or the C-stick across) to roll, ZR and ZL (X and
// B, or the C-stick up and down) for speed, R or A to fire.
func controlsOf(keys *engine.Input) (flight.Controls, bool) {
	stickX, stickY := keys.Stick()
	cStickX, cStickY := keys.CStick()
	controls := flight.Controls{
		Pitch:    -stickY,
		Yaw:      stickX,
		Roll:     keys.Axis(engine.KeyLeft, engine.KeyRight) + cStickX,
		Throttle: keys.Axis(engine.ButtonZL, engine.ButtonZR) + keys.Axis(engine.KeyEscape, engine.KeyX) + cStickY,
	}
	return controls, keys.Down(engine.ButtonR) || keys.Down(engine.KeySpace)
}

// pausePressed tells whether START was pressed this frame: it opens and closes the pause menu, as
// in the 3DS program.
func pausePressed(keys *engine.Input) bool { return keys.Pressed(engine.ButtonStart) }
