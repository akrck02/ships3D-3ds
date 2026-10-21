package main

import (
	"math"

	"github.com/akrck02/go-engine/engine"
	"github.com/akrck02/go-engine/space"
	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/models"

	"github.com/jonbul/ships3D-3ds/arena"
	"github.com/jonbul/ships3D-3ds/chase"
	"github.com/jonbul/ships3D-3ds/match"
	"github.com/jonbul/ships3D-3ds/protocol"
	"github.com/jonbul/ships3D-3ds/sounds"
	"github.com/jonbul/ships3D-3ds/transport"
)

// engineFar is how far go-engine's camera sees (v0.0.3 fixes it): the sky is brought inside it.
const engineFar = 2000

// spaceAmbient is the light that reaches the ships from all round, with no sky to light them.
const spaceAmbient = 0.9

// The volumes, from 0 to 1, until the game has its settings.
const (
	effectsVolume = 0.8
	musicVolume   = 0.5
)

// Pilot plays a match on the server: it joins with a ship, flies it with the player's keys, and
// shows the arena (package arena) inside the engine's scene, its camera chasing the ship.
type Pilot struct {
	Server string `json:"server"` // address:port, or https://name:port; empty is the real server
	Ship   string `json:"ship"`   // the ship's id; empty is the Arrow
	Name   string `json:"name"`   // a guest's name

	world  *arena.World
	sound  *sounds.Sounds
	played *match.Match
	view   models.Camera
	status string
	// readings is the HUD's drawing, made again every frame (readings*.go).
	readings draw2d.List
}

// The real server, and what a guest flies when the scene names nothing.
const (
	realServer  = "https://jonbul.ddns.net:3000"
	defaultShip = "default:arrow"
	guestName   = "Pilot"
)

// Start builds the arena in the engine's scene and starts the match.
func (pilot *Pilot) Start(context *engine.Context) {
	pilot.world = arena.NewWorldIn(context.Scene.World())
	pilot.world.FitSky(engineFar)
	// The desktop paints the sky as a blend from its zenith to its horizon, and the zenith is the
	// sky's light on the ships: space is dark all round, so the ships take their light from the
	// ambient instead.
	environment := &context.Scene.World().Environment
	environment.SkyZenith = environment.SkyHorizon
	environment.Ambient = spaceAmbient
	pilot.view = models.Camera{Far: engineFar, Near: arena.CameraNear, FieldOfView: chase.FieldOfView}
	pilot.startSound()
	played, err := match.Start(transport.Dial, orDefault(pilot.Server, realServer), "",
		orDefault(pilot.Ship, defaultShip), orDefault(pilot.Name, guestName), match.Hooks{
			Fired: func() { pilot.sound.Play(pilot.sound.Laser, 1) },
			Received: func(message protocol.Incoming, _ int) {
				if message.Type == protocol.TypeFire {
					pilot.sound.Play(pilot.sound.Laser, sounds.OtherShotsVolume)
				}
			},
		})
	if err != nil {
		pilot.status = err.Error()
		println(pilot.status) // the debug log: the arena shows empty
		return
	}
	pilot.played = played
}

// Update moves the match on a frame and shows it, through the scene's camera.
func (pilot *Pilot) Update(context *engine.Context) {
	if pilot.played == nil {
		return
	}
	controls, firing := controlsOf(context.Input)
	if err := pilot.played.Exchange(controls, firing, context.Delta, context.Time); err != nil {
		pilot.status = err.Error()
		println(pilot.status) // the debug log
		pilot.played.Close()
		pilot.played = nil
		return
	}
	state := pilot.played.Game
	events := state.TakeEvents()
	pilot.sound.Hear(events, state)
	flying := state.Welcomed && state.Alive
	pilot.sound.RunEngine(flying, math.Abs(state.Own.Speed)/math.Max(state.Settings.MaxSpeed, 1))
	pilot.world.Show(state, events, &pilot.view, context.Delta, context.Time)
	pilot.placeCamera(context)
	showReadings(context.Scene, &pilot.readings, state, context.Time)
}

// startSound opens the sound and starts the music; without it the game plays on in silence.
func (pilot *Pilot) startSound() {
	sound, err := sounds.New()
	if err != nil {
		println("no sound: " + err.Error()) // the debug log
	}
	pilot.sound = sound
	pilot.sound.SetVolumes(effectsVolume, musicVolume)
}

// placeCamera puts the scene's main camera where the arena's chase camera is.
func (pilot *Pilot) placeCamera(context *engine.Context) {
	camera, found := context.Scene.MainCamera()
	if !found {
		return
	}
	eye := vectorOf(pilot.view.Position[0], pilot.view.Position[1], pilot.view.Position[2])
	target := vectorOf(pilot.view.Target[0], pilot.view.Target[1], pilot.view.Target[2])
	up := vectorOf(pilot.view.Up[0], pilot.view.Up[1], pilot.view.Up[2])
	camera.Transform.Position = eye
	camera.Transform.Rotation = space.RotationLookingAt(target.Subtract(eye), up)
}

func vectorOf(x, y, z float32) space.Vector {
	return space.NewVector(float64(x), float64(y), float64(z))
}

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
