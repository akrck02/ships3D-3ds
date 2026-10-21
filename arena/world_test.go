package arena

import (
	"testing"

	"github.com/akrck02/go-renderer/models"

	"github.com/jonbul/ships3D-3ds/chase"
	"github.com/jonbul/ships3D-3ds/game"
)

func TestTheMatchSeesAsWideAsTheChaseCameraAfterTheTitle(test *testing.T) {
	world := NewWorld()
	var view models.Camera
	world.ShowTitle(0, &view)
	world.Show(game.New(), nil, &view, 0, 0)
	if view.FieldOfView != chase.FieldOfView {
		test.Errorf("the match sees %v radians across, want the chase camera's %v", view.FieldOfView, chase.FieldOfView)
	}
}
