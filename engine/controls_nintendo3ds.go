//go:build nintendo3ds

package main

import (
	"github.com/akrck02/go-engine/engine"
	"github.com/akrck02/go-renderer/ctr"

	"github.com/jonbul/ships3D-3ds/flight"
)

// circleReach is how far the circle pad goes, in its own steps.
const circleReach = 156

// controlsOf reads the console as the 3DS program does: the circle pad for the nose (down raises
// it) and the turn, the cross's left and right to roll, ZR and ZL (or X and B) for speed, R or A
// to fire. go-engine gives a game only some of the buttons as keys, so they are read here.
func controlsOf(*engine.Input) (flight.Controls, bool) {
	reading := ctr.ReadControls()
	held := reading.Held
	controls := flight.Controls{
		Pitch: -float64(reading.CircleY) / circleReach,
		Yaw:   float64(reading.CircleX) / circleReach,
	}
	if held&ctr.ButtonLeft != 0 {
		controls.Roll--
	}
	if held&ctr.ButtonRight != 0 {
		controls.Roll++
	}
	if held&(ctr.ButtonZR|ctr.ButtonX) != 0 {
		controls.Throttle++
	}
	if held&(ctr.ButtonZL|ctr.ButtonB) != 0 {
		controls.Throttle--
	}
	return controls, held&(ctr.ButtonR|ctr.ButtonA) != 0
}
