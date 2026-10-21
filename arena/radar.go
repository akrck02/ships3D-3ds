package arena

import (
	"math"

	"github.com/akrck02/go-renderer/graphics"
	"github.com/akrck02/go-renderer/models"
	"github.com/akrck02/go-renderer/scene"

	"github.com/jonbul/ships3D-3ds/flight"
	"github.com/jonbul/ships3D-3ds/protocol"
	"github.com/jonbul/ships3D-3ds/radar"
	"github.com/jonbul/ships3D-3ds/ship"
)

// The 3D radar: a dish of radius 1, seen from above and behind our ship as the browser's radar is,
// with a blip for each ship the sweep last saw and a stalk from the dish's plane to it.
const (
	radarBlips        = 64 // most blips drawn at once
	radarBlipSize     = 0.09
	radarStalkWidth   = 0.018
	radarRingWidth    = 0.012
	radarRingSegments = 48
	radarSweepWidth   = 0.02
	radarMarkerSize   = 0.16
	radarUnseen       = 0.08 // a blip dimmer than this isn't drawn
)

// The radar's colours, linear as go-renderer takes them: the flat radar's, from its sRGB hex.
var (
	radarRingColour  = [3]float32{0.047, 0.13, 0.18}   // #3d6577
	radarFaintColour = [3]float32{0.018, 0.053, 0.078} // #24414f
	radarSweepColour = [3]float32{0.078, 0.546, 0.93}  // #4fc3f7
	radarEnemyColour = graphics.Vec4{1, 0.085, 0.085, 1}
	radarFarColour   = graphics.Vec4{0.254, 0.034, 0.034, 1}
	// RadarBackground is the radar's dark, what its screen is cleared to.
	RadarBackground = graphics.Vec4{0.0044, 0.010, 0.016, 1} // #0e1a22, the touch screen's dark
)

// RadarView looks at the dish from above and behind, as the browser's radar camera does.
var RadarView = models.Camera{
	Position: graphics.Vec4{0, 1.55, -2.25, 1},
	Target:   graphics.Vec4{0, -0.1, 0.18, 1},
	Up:       graphics.Vec4{0, 1, 0, 0},
	Near:     0.1,
	Far:      20,
}

// Radar is the 3D radar's scene and what moves in it.
type Radar struct {
	Scene  *scene.Scene
	blips  []*scene.Node // each with a mesh of its own, so each can fade on its own
	stalks []*scene.Node
	sweep  *scene.Node
}

// NewRadar builds the radar in a scene of its own.
func NewRadar() *Radar { return NewRadarIn(scene.New()) }

// NewRadarIn builds the radar inside a scene that already exists, as a game engine's.
func NewRadarIn(dish *scene.Scene) *Radar {
	// By day, so nothing is dimmed for the night; the sun straight up, so there is no twilight
	// glow on the background.
	dish.Environment.SunDirection = graphics.Vec4{0, 1, 0, 0}
	dish.Environment.SkyZenith, dish.Environment.SkyHorizon = RadarBackground, RadarBackground
	dish.Environment.FogColor = RadarBackground
	still := func(name string, mesh *scene.Mesh) {
		dish.Nodes = append(dish.Nodes, &scene.Node{Name: name, Transform: graphics.Identity(), Mesh: mesh})
	}
	still("edge", ring(1, radarRingColour))
	still("half", ring(0.5, radarFaintColour))
	still("cross", crossLines())
	built := &Radar{Scene: dish}
	built.sweep = &scene.Node{Name: "sweep", Mesh: Box("sweep", radarSweepWidth, radarSweepWidth, 1, radarSweepColour), Moving: true}
	// The sweep's box runs from the centre forward: moved half its length, then turned.
	built.sweep.Mesh.Positions = shifted(built.sweep.Mesh.Positions, 0, 0, 0.5)
	dish.Nodes = append(dish.Nodes, built.sweep, &scene.Node{Name: "own", Transform: graphics.Identity(), Mesh: ownMarker()})
	blip := solid(ship.Build([]protocol.Layer{{Visible: true, Parts: []protocol.Part{{
		Type: "sphere", Scale: flight.Vector{radarBlipSize, radarBlipSize, radarBlipSize}, Color: "#ffffff"}}}}, ship.Coarse))
	stalk := Box("stalk", radarStalkWidth, 1, radarStalkWidth, [3]float32{1, 1, 1})
	for range radarBlips {
		blipNode := &scene.Node{Name: "blip", Mesh: WithOwnMaterial(blip), Hidden: true, Moving: true}
		stalkNode := &scene.Node{Name: "stalk", Mesh: WithOwnMaterial(stalk), Hidden: true, Moving: true}
		built.blips, built.stalks = append(built.blips, blipNode), append(built.stalks, stalkNode)
		dish.Nodes = append(dish.Nodes, blipNode, stalkNode)
	}
	return built
}

// Show places the sweep and the blips for what it has seen by now.
func (shown3D *Radar) Show(scanner *radar.Scanner, now float64) {
	angle := radar.Angle(now)
	shown3D.sweep.Transform = graphics.RotateY(-angle)
	drawn := 0
	for _, echo := range scanner.Echoes {
		if drawn == len(shown3D.blips) {
			break
		}
		x, y, z := echo.Blip.Scene()
		colour := radarEnemyColour
		if !echo.Blip.InRange {
			colour = radarFarColour
		}
		brightness := float32(echo.Brightness(now))
		if brightness < radarUnseen {
			continue // faded out: on the dark dish it would only be a dark speck
		}
		faded := graphics.Vec4{colour[0] * brightness, colour[1] * brightness, colour[2] * brightness, 1}
		blip, stalk := shown3D.blips[drawn], shown3D.stalks[drawn]
		blip.Hidden, stalk.Hidden = false, math.Abs(y) < radarBlipSize/2
		blip.Transform = graphics.Translate(graphics.Vec4{float32(x), float32(y), float32(z), 0})
		stalk.Transform = graphics.Translate(graphics.Vec4{float32(x), float32(y / 2), float32(z), 0}).
			Multiply(graphics.Scale(graphics.Vec4{1, float32(math.Abs(y)), 1, 0}))
		blip.Mesh.Material.BaseColor, stalk.Mesh.Material.BaseColor = faded, faded
		drawn++
	}
	for index := drawn; index < len(shown3D.blips); index++ {
		shown3D.blips[index].Hidden, shown3D.stalks[index].Hidden = true, true
	}
}

// solid is a mesh drawn in its own flat colours, without light: the radar is a screen, not a
// lit object.
func solid(mesh *scene.Mesh) *scene.Mesh {
	mesh.Material = &scene.Material{Name: "radar", Kind: scene.KindUnlit, BaseColor: graphics.Vec4{1, 1, 1, 1}}
	return mesh
}

// ring is a flat circle of radius in the dish's plane, a band radarRingWidth wide.
func ring(radius float64, colour [3]float32) *scene.Mesh {
	mesh := solid(&scene.Mesh{Name: "ring"})
	inner, outer := float32(radius-radarRingWidth/2), float32(radius+radarRingWidth/2)
	for segment := range radarRingSegments {
		first := float64(segment) / radarRingSegments * 2 * math.Pi
		second := float64(segment+1) / radarRingSegments * 2 * math.Pi
		corners := [4][3]float32{
			{inner * float32(math.Cos(first)), 0, inner * float32(math.Sin(first))},
			{outer * float32(math.Cos(first)), 0, outer * float32(math.Sin(first))},
			{outer * float32(math.Cos(second)), 0, outer * float32(math.Sin(second))},
			{inner * float32(math.Cos(second)), 0, inner * float32(math.Sin(second))},
		}
		for _, index := range [6]int{0, 1, 2, 0, 2, 3} {
			mesh.Positions = append(mesh.Positions, corners[index][0], corners[index][1], corners[index][2])
			mesh.Normals = append(mesh.Normals, 0, 1, 0)
			mesh.Colors = append(mesh.Colors, colour[0], colour[1], colour[2], 1)
		}
	}
	return mesh
}

// crossLines are the dish's two diameters: ahead to behind, left to right.
func crossLines() *scene.Mesh {
	mesh := solid(&scene.Mesh{Name: "cross"})
	AppendBox(mesh, [3]float32{}, radarRingWidth, radarRingWidth, 2, radarFaintColour)
	AppendBox(mesh, [3]float32{}, 2, radarRingWidth, radarRingWidth, radarFaintColour)
	return mesh
}

// ownMarker is a small cone at the centre, its tip ahead: our ship.
func ownMarker() *scene.Mesh {
	marker := ship.Build([]protocol.Layer{{Visible: true, Parts: []protocol.Part{{
		Type: "cone", Rotation: flight.Vector{math.Pi / 2, 0, 0},
		Scale: flight.Vector{radarMarkerSize * 0.6, radarMarkerSize, radarMarkerSize * 0.4}, Color: "#4fc3f7"}}}}, ship.Coarse)
	return solid(marker)
}

// shifted moves every position of a mesh's xyz triples.
func shifted(positions []float32, x, y, z float32) []float32 {
	for index := 0; index+2 < len(positions); index += 3 {
		positions[index] += x
		positions[index+1] += y
		positions[index+2] += z
	}
	return positions
}
