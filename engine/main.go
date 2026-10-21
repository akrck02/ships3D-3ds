// Ships 3D, made with go-engine: the same game as the browser's and the 3DS program's, built for
// the desktop and the 3DS from one project. The editor keeps the scenes in scenes/; the game's
// code is the Go beside this file and the repository's packages (arena, match, game…). Play it
// with `go run ./engine`.
package main

import (
	"embed"
	"log"

	"github.com/akrck02/go-engine/engine"
)

// files are the project's settings, scenes and assets, carried inside the game.
//
//go:embed engine.json scenes all:assets
var files embed.FS

func main() {
	// Each behaviour the scenes use is registered here, by its type's name.
	engine.Register(Pilot{})
	if err := engine.Run(files); err != nil {
		log.Fatal(err)
	}
}
