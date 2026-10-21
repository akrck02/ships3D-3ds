package main

import (
	"github.com/akrck02/go-engine/engine"
	"github.com/akrck02/go-renderer/graphics"
	"github.com/akrck02/go-renderer/scene"

	"github.com/jonbul/ships3D-3ds/arena"
	"github.com/jonbul/ships3D-3ds/flight"
	"github.com/jonbul/ships3D-3ds/protocol"
	"github.com/jonbul/ships3D-3ds/ship"
)

// backdropRadius is the black sphere's: past the stars (arena.SkyRadius), inside the camera's
// reach (arena.CameraFar).
const backdropRadius = (arena.SkyRadius + arena.CameraFar) / 2

// backdrop is space's black, all round the camera. The desktop paints its own sky (a day's, with
// the sun's glow) wherever nothing covers the far plane; this sphere, behind the stars, covers it.
// The 3DS paints no sky, only its clear colour, so there it changes nothing.
func backdrop(context *engine.Context) *scene.Node {
	mesh := ship.Build([]protocol.Layer{{Visible: true, Parts: []protocol.Part{{
		Type: "sphere", Scale: flight.Vector{2 * backdropRadius, 2 * backdropRadius, 2 * backdropRadius}, Color: "#02030a"}}}}, ship.Coarse)
	mesh.Material = &scene.Material{Name: "backdrop", Kind: scene.KindUnlit, BaseColor: graphics.Vec4{1, 1, 1, 1}, DoubleSided: true}
	node := &scene.Node{Name: "backdrop", Mesh: mesh, Transform: graphics.Identity(), Moving: true}
	world := context.Scene.World()
	world.Nodes = append(world.Nodes, node)
	return node
}

// shownBackdrop is the backdrop of the scene playing now.
var shownBackdrop *scene.Node

// followCamera keeps the backdrop round the camera.
func followCamera(played *engine.Scene) {
	if shownBackdrop == nil {
		return
	}
	if camera, found := played.MainCamera(); found {
		position := camera.Transform.Position
		shownBackdrop.Transform = graphics.Translate(graphics.Vec4{float32(position.X), float32(position.Y), float32(position.Z), 0})
	}
}
