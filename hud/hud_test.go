package hud

import (
	"testing"

	"github.com/akrck02/go-renderer/draw2d"

	"github.com/jonbul/ships3D-3ds/flight"
	"github.com/jonbul/ships3D-3ds/game"
)

func TestTheLifeBarIsDrawnAndTheCopilotOnlyNearTheEdge(test *testing.T) {
	state := game.New()
	state.Settings.MaxLife, state.Settings.WorldRadius = 10, 1500
	state.Player.Life = 5
	var middle, edge draw2d.List
	Paint(&middle, state, 0, 400, 240)
	if middle.IsEmpty() {
		test.Fatal("nothing drawn: the life bar is missing")
	}
	state.Own.Position = flight.Vector{0, 0, 1490}
	Paint(&edge, state, 0, 400, 240)
	if len(edge.Commands) <= len(middle.Commands) {
		test.Fatalf("%d commands at the edge, %d in the middle: the copilot is missing", len(edge.Commands), len(middle.Commands))
	}
}
