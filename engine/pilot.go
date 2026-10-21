package main

import (
	"math"

	"github.com/akrck02/go-engine/engine"
	"github.com/akrck02/go-engine/space"
	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/models"

	"github.com/jonbul/ships3D-3ds/arena"
	"github.com/jonbul/ships3D-3ds/chase"
	"github.com/jonbul/ships3D-3ds/flight"
	"github.com/jonbul/ships3D-3ds/hud"
	"github.com/jonbul/ships3D-3ds/match"
	"github.com/jonbul/ships3D-3ds/pixelart"
	"github.com/jonbul/ships3D-3ds/protocol"
	"github.com/jonbul/ships3D-3ds/sounds"
	"github.com/jonbul/ships3D-3ds/transport"
)

// spaceAmbient is the light that reaches the ships from all round, with no sky to light them.
const spaceAmbient = 0.9

// screenPixelsHigh is the 3DS's top screen's height: the HUD is laid out in its pixels.
const screenPixelsHigh = 240

// Pilot plays a match on the server: it joins with a ship, flies it with the player's keys, and
// shows the arena (package arena) inside the engine's scene, its camera chasing the ship.
type Pilot struct {
	Server string `json:"server"` // address:port, or https://name:port; empty is the real server

	world  *arena.World
	sound  *soundSet
	played *match.Match
	view   models.Camera
	status string
	now    float64 // seconds, this frame
	paused bool    // the pause menu is open
	resume bool    // which of the pause menu's two is chosen: Resume, or Leave
}

// The real server, and what a guest flies when the scene names nothing.
const (
	realServer  = "https://jonbul.ddns.net:3000"
	defaultShip = "default:arrow"
	guestName   = "Pilot"
)

// Start builds the arena in the engine's scene and starts the match.
func (pilot *Pilot) Start(context *engine.Context) {
	pilot.world = stage(context)
	drawOver(pilot.paintReadings)
	pilot.view = models.Camera{Far: arena.CameraFar, Near: arena.CameraNear, FieldOfView: chase.FieldOfView}
	pilot.sound = gameSound
	played, err := match.Start(transport.Dial, orDefault(pilot.Server, realServer), "",
		orDefault(chosen.shipID, defaultShip), orDefault(chosen.name, guestName), match.Hooks{
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
	if pausePressed(context.Input) {
		pilot.paused, pilot.resume = !pilot.paused, true
		playSound(func(sound *soundSet) { sound.Play(sound.Choose, 1) })
	}
	if pilot.paused && pilot.choosePaused(context) {
		return
	}
	controls, firing := controlsOf(context.Input)
	if pilot.paused {
		controls, firing = flight.Controls{}, false // the ship flies on straight while paused
	}
	if err := pilot.played.Exchange(controls, firing, context.Delta, context.Time); err != nil {
		pilot.status = err.Error()
		println(pilot.status) // the debug log
		pilot.played.Close()
		pilot.played = nil
		return
	}
	pilot.now = context.Time
	darkenSky(context.Scene)
	state := pilot.played.Game
	events := state.TakeEvents()
	pilot.sound.Hear(events, state)
	flying := state.Welcomed && state.Alive
	pilot.sound.RunEngine(flying, math.Abs(state.Own.Speed)/math.Max(state.Settings.MaxSpeed, 1))
	pilot.world.Show(state, events, &pilot.view, context.Delta, context.Time)
	pilot.placeCamera(context)
}

// choosePaused moves the pause menu's choice and acts on it; it tells whether the match was left.
func (pilot *Pilot) choosePaused(context *engine.Context) bool {
	keys := context.Input
	if keys.Pressed(engine.KeyUp) || keys.Pressed(engine.KeyDown) {
		pilot.resume = !pilot.resume
		playSound(func(sound *soundSet) { sound.Play(sound.Move, 1) })
	}
	if !confirmed(keys) {
		return false
	}
	if pilot.resume {
		pilot.paused = false
		playSound(func(sound *soundSet) { sound.Play(sound.Back, 1) })
		return false
	}
	pilot.played.Close()
	pilot.played = nil
	playSound(func(sound *soundSet) {
		sound.Play(sound.Back, 1)
		sound.RunEngine(false, 0)
	})
	context.ChangeScene(shipScene)
	return true
}

// paintReadings paints the HUD over the 3D: the 3DS's life bar and the copilot.
func (pilot *Pilot) paintReadings(drawing *draw2d.List, width, height float32) {
	if pilot.played == nil {
		return
	}
	// Laid out in the 3DS's pixels, and made larger on a larger screen, so it looks the same.
	scale := float32(max(1, math.Floor(float64(height)/screenPixelsHigh)))
	start := len(drawing.Commands)
	hud.Paint(drawing, pilot.played.Game, pilot.now, width/scale, height/scale)
	if pilot.paused {
		pilot.paintPause(drawing, width/scale, height/scale)
	}
	pixelart.Scale(drawing, start, scale)
}

// The pause menu's measures, in pixels.
const (
	pauseWidth  = 160
	pauseHeight = 96
	pauseRow    = 22
)

// paintPause draws the pause menu: Resume, and Leave for the ships.
func (pilot *Pilot) paintPause(drawing *draw2d.List, width, height float32) {
	area := draw2d.Rect{X: (width - pauseWidth) / 2, Y: (height - pauseHeight) / 2, Width: pauseWidth, Height: pauseHeight}
	window(drawing, area, "PAUSED")
	for index, label := range []string{"Resume", "Leave"} {
		top := area.Y + 32 + float32(index)*pauseRow
		if (index == 0) == pilot.resume {
			pixelart.Stepped(drawing, draw2d.Rect{X: area.X + 8, Y: top, Width: area.Width - 16, Height: pauseRow - 4}, 1, colours.chosen)
			cursor(drawing, area.X+14, top+5, pilot.now)
		}
		drawing.Text(label, area.X+32, top+14, pixelart.TextSize, draw2d.Regular, colours.ink)
	}
}

// darkenSky keeps the sky the dark of space: the desktop paints it as a blend from its zenith to
// its horizon, and the zenith is the sky's light on the ships, so the ships take their light from
// the ambient instead. The engine sets its own every frame, so this follows it.
func darkenSky(scene *engine.Scene) {
	environment := &scene.World().Environment
	environment.SkyZenith = environment.SkyHorizon
	environment.Ambient = spaceAmbient
	followCamera(scene)
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
