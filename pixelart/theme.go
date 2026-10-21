package pixelart

import (
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// Theme is a go-ui theme in the game's look, over base: night-blue glass, starlight ink, a sky
// blue for what is chosen, the pixel fonts' sizes, square corners and hard shadows.
func Theme(base ui.Theme) ui.Theme {
	theme := base
	theme.Background = draw2d.Hex("#0a1430")
	theme.Surface = draw2d.Hex("#12224a")
	theme.Sunken = draw2d.Hex("#1a2d5c")
	theme.Line = draw2d.Hex("#3a5590")
	theme.ControlLine = draw2d.Hex("#9cc8f2")
	theme.Ink = draw2d.Hex("#eaf4ff")
	theme.InkSoft = draw2d.Hex("#a9bedd")
	theme.InkQuiet = draw2d.Hex("#7f93b8")
	theme.Primary = draw2d.Hex("#5fb0ff")
	theme.OnPrimary = draw2d.Hex("#02060f")
	theme.PrimaryTint = draw2d.Hex("#2a4a80")
	theme.FocusRing = draw2d.Hex("#5fb0ff80")
	theme.Shadow = draw2d.Hex("#02060f")
	theme.TextSize = TextSize
	theme.SmallTextSize = TextSize
	theme.HeadingSize = TextSize
	theme.ControlRadius = 0
	theme.PanelRadius = 0
	theme.ShadowBlur = 0 // a hard shadow, a pixel's step below
	theme.ShadowDrop = 2
	return theme
}
