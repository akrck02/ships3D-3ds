//go:build nintendo3ds

package main

import (
	"github.com/akrck02/go-3ds/bottom"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"

	"github.com/jonbul/ships3D-3ds/pixelart"
)

// The game's look on the console: pixel fonts (package pixelart) and square corners, as handheld
// space shooters had, in the colours of space.

// The pixel fonts' sizes, in pixels.
const (
	textSize    = 16 // VT323: text, captions
	headingSize = 16 // Silkscreen: headings, buttons' labels in bold
	titleSize   = 40 // Silkscreen: the game's name
)

// usePixelFonts makes every text draw in the pixel fonts; call it before the first painter is
// made, which reads them.
func usePixelFonts() { pixelart.UseFonts() }

// spaceTheme is go-3ds's bottom-screen theme in pixel art and the colours of space: a night blue,
// pale ink, and an arcade gold for what is chosen; square corners and hard shadows.
func spaceTheme() ui.Theme {
	theme := bottom.Theme()
	theme.Background = draw2d.Hex("#070b1a")
	theme.Surface = draw2d.Hex("#0e1530")
	theme.Sunken = draw2d.Hex("#18214a")
	theme.Line = draw2d.Hex("#2b3a78")
	theme.ControlLine = draw2d.Hex("#4458b0")
	theme.Ink = draw2d.Hex("#e9edff")
	theme.InkSoft = draw2d.Hex("#aab4ea")
	theme.InkQuiet = draw2d.Hex("#7a86c4")
	theme.Primary = draw2d.Hex("#ffcc33")
	theme.OnPrimary = draw2d.Hex("#1c1400")
	theme.PrimaryTint = draw2d.Hex("#3a3112")
	theme.FocusRing = draw2d.Hex("#ffcc3380")
	theme.Shadow = draw2d.Hex("#000000")
	theme.TextSize = textSize
	theme.SmallTextSize = textSize
	theme.HeadingSize = headingSize
	theme.IconSize = 16
	theme.ControlRadius = 0
	theme.PanelRadius = 0
	theme.ShadowBlur = 0 // a hard shadow, a pixel's step below
	theme.ShadowDrop = 2
	return theme
}
