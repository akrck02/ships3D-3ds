//go:build !nintendo3ds

package main

import (
	"github.com/akrck02/go-engine/engine"

	"github.com/jonbul/ships3D-3ds/flight"
	"github.com/jonbul/ships3D-3ds/profile"
)

// controlsOf reads the keys as the flight model's controls, as the browser does: the arrows (or
// a gamepad's stick) for the nose and the turn, A and D to roll, W and S for speed, Space to fire.
// Fine aim softens a stick near its centre; inverting the pitch swaps the nose's keys, as on the
// 3DS.
func controlsOf(keys *engine.Input, settings profile.Settings) (flight.Controls, bool) {
	turn, nose := keys.Axis(engine.KeyLeft, engine.KeyRight), keys.Axis(engine.KeyDown, engine.KeyUp)
	if stick := keys.State(); stick != nil {
		precision := float64(settings.Precision)
		turn, nose = turn+flight.Curve(stick.StickX, precision), nose-flight.Curve(stick.StickY, precision)
	}
	if settings.InvertPitch {
		nose = -nose
	}
	return flight.Controls{
		Pitch:    nose,
		Yaw:      turn,
		Roll:     keys.Axis(engine.KeyA, engine.KeyD),
		Throttle: keys.Axis(engine.KeyS, engine.KeyW),
	}, keys.Down(engine.KeySpace)
}

// turnTilt is the console's tilt, which a computer has none of.
func turnTilt(*engine.Input, float64) {}

// menuPressed tells whether the player opened or closed the menu: Tab or P, the keyboard's START
// (Enter launches and fires).
func menuPressed(keys *engine.Input) bool {
	return keys.Pressed(engine.KeyTab) || keys.Pressed(engine.KeyP)
}
