package main

import (
	"github.com/akrck02/go-engine/engine"
	"github.com/akrck02/go-engine/scenefile"

	"github.com/jonbul/ships3D-3ds/arena"
	"github.com/jonbul/ships3D-3ds/game"
	"github.com/jonbul/ships3D-3ds/radar"
)

// radarScreen is the 3D radar on the second screen (the 3DS's bottom one), as in the 3DS program:
// a scene under the interface, its camera above and behind the dish. The desktop has no second
// screen: the radar still plays there, and the editor's Play shows it under the game with -3ds.
//
// TODO(akrck02): show it on the desktop too, in a corner as the browser does; that needs a
// second view in go-engine's window.
type radarScreen struct {
	scene   *engine.Scene
	camera  *engine.Node
	dish    *arena.Radar
	scanner *radar.Scanner
	blips   map[string]radar.Blip
}

// openRadar builds the radar's scene and puts it on the second screen.
func openRadar() *radarScreen {
	bottom := engine.NewScene(files)
	shown := &radarScreen{scene: bottom, dish: arena.NewRadarIn(bottom.World()), scanner: radar.NewScanner(), blips: map[string]radar.Blip{}}
	camera, err := bottom.NewNode(scenefile.Node{Name: "Camera", Camera: &scenefile.Camera{Main: true, FieldOfView: radarFieldOfView, Near: 0.1, Far: 20}})
	if err != nil {
		println(err.Error()) // the debug log: the second screen stays empty
		return shown
	}
	bottom.Root.Add(camera)
	shown.camera = camera
	engine.SetSecondScreen(bottom)
	return shown
}

// radarFieldOfView is the radar camera's, in degrees, as go-renderer's default.
const radarFieldOfView = 45

// show turns the sweep and lets it see the ships it passes over.
func (shown *radarScreen) show(state *game.Game, now float64) {
	clear(shown.blips)
	for id, remote := range state.Remotes {
		if !remote.Player.Dead {
			shown.blips[id] = radar.Place(remote.Shown.Position, state.Own.Position, state.Own.Rotation, radar.Ranges[radar.DefaultRange])
		}
	}
	shown.scanner.Update(now, shown.blips)
	shown.dish.Show(shown.scanner, now)
	if shown.camera != nil {
		view := arena.RadarView
		eye := vectorOf(view.Position[0], view.Position[1], view.Position[2])
		target := vectorOf(view.Target[0], view.Target[1], view.Target[2])
		shown.camera.Transform.Position = eye
		shown.camera.Transform.LookAt(target)
	}
}

// close takes the radar off the second screen, as the match ends.
func (*radarScreen) close() { engine.SetSecondScreen(nil) }
