//go:build nintendo3ds

package main

import (
	"github.com/akrck02/go-engine/engine"

	"github.com/jonbul/ships3D-3ds/flight"
	"github.com/jonbul/ships3D-3ds/profile"
	"github.com/jonbul/ships3D-3ds/tilt"
)

// controlsOf reads the console as the 3DS program does: the circle pad for the nose (down raises
// it, or up when the settings say so) and the turn, softer near its centre with fine aim; the
// console's tilt too when steering by tilting; the cross's left and right (or the C-stick across)
// to roll, ZR and ZL (X and B, or the C-stick up and down) for speed, R or A to fire.
func controlsOf(keys *engine.Input, settings profile.Settings) (flight.Controls, bool) {
	stickX, stickY := keys.Stick()
	rightStickX, rightStickY := keys.RightStick()
	precision := float64(settings.Precision)
	pitch := -flight.Curve(stickY, precision)
	if settings.InvertPitch {
		pitch = -pitch
	}
	controls := flight.Controls{
		Pitch:    pitch,
		Yaw:      flight.Curve(stickX, precision),
		Roll:     keys.Axis(engine.KeyLeft, engine.KeyRight) + rightStickX,
		Throttle: keys.Axis(engine.ButtonLeftTrigger, engine.ButtonRightTrigger) + keys.Axis(engine.KeyEscape, engine.KeyX) + rightStickY,
	}
	if settings.Gyroscope {
		controls.Pitch += tilt.Axis(tilted.Pitch)
		controls.Yaw += tilt.Axis(tilted.Steer)
	}
	return controls, keys.Down(engine.ButtonRightShoulder) || keys.Down(engine.KeySpace)
}

// turnTilt adds this frame's turning of the console up, and L centres it, as in the 3DS program.
// The gyroscope's roll tips the top screen (toward the player is negative); its pitch turns the
// console like a wheel.
func turnTilt(keys *engine.Input, seconds float64) {
	turning = keys.Gyroscope()
	tilted.Turn(-turning.Roll, turning.Pitch, seconds)
	if keys.Pressed(engine.ButtonLeftShoulder) {
		tilted.Centre()
	}
}

// menuPressed tells whether START was pressed this frame: it opens and closes the menu, as in
// the 3DS program.
func menuPressed(keys *engine.Input) bool { return keys.Pressed(engine.ButtonStart) }
