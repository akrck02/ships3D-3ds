package pixelart

import (
	_ "embed"

	"github.com/akrck02/go-renderer/opengl/overlay"
)

// The game's pixel fonts, under the SIL Open Font License (fonts/OFL-*.txt), from Google Fonts:
// VT323 for text, Silkscreen for headings and names, Press Start 2P for figures. Each is crisp at
// its own sizes only: VT323 at 16, the others at multiples of 8.
var (
	//go:embed fonts/VT323-Regular.ttf
	vt323 []byte
	//go:embed fonts/Silkscreen-Bold.ttf
	silkscreenBold []byte
	//go:embed fonts/PressStart2P-Regular.ttf
	pressStart []byte
)

// UseFonts makes every text draw in the pixel fonts: VT323 as the regular and medium weight,
// Silkscreen as the bold, Press Start 2P as the mono. Call it before the first painter is made,
// which reads them.
func UseFonts() {
	overlay.UseFonts(overlay.Fonts{Regular: vt323, Medium: vt323, Bold: silkscreenBold, Mono: pressStart})
}
