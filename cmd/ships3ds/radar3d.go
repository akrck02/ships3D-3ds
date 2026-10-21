//go:build nintendo3ds

package main

import (
	"github.com/akrck02/go-renderer/ctr"

	"github.com/jonbul/ships3D-3ds/arena"
	"github.com/jonbul/ships3D-3ds/radar"
)

// radarScene is the 3D radar (arena.Radar) and the renderer that draws it on the touch screen.
type radarScene struct {
	radar    *arena.Radar
	renderer *ctr.SceneRenderer
}

func newRadarScene(screen *ctr.Screen) (*radarScene, error) {
	renderer, err := ctr.NewSceneRenderer(screen)
	if err != nil {
		return nil, err
	}
	return &radarScene{radar: arena.NewRadar(), renderer: renderer}, nil
}

// paint draws the radar on the bottom screen, in a frame begun with ctr.BeginFrame, for what the
// sweep has seen by now.
func (shown3D *radarScene) paint(surface ctr.Surface, scanner *radar.Scanner, now float64) {
	shown3D.radar.Show(scanner, now)
	shown3D.renderer.PaintOn(surface, shown3D.radar.Scene, arena.RadarView)
}
