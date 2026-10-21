package main

import (
	"github.com/akrck02/go-engine/engine"
	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/models"

	"github.com/jonbul/ships3D-3ds/arena"
	"github.com/jonbul/ships3D-3ds/pixelart"
	"github.com/jonbul/ships3D-3ds/touchscreen"
)

// The scenes, in the order a player goes through them.
const (
	titleScene = "scenes/title.scene.json"
	modeScene  = "scenes/mode.scene.json"
	shipScene  = "scenes/ships.scene.json"
	arenaScene = "scenes/arena.scene.json"
)

// chosen is what the player has picked, carried from one scene to the next.
var chosen = struct{ shipID string }{shipID: defaultShip}

// confirmed tells whether the player said yes this frame: Enter, or Space (the 3DS's A).
func confirmed(keys *engine.Input) bool {
	return keys.Pressed(engine.KeyEnter) || keys.Pressed(engine.KeySpace)
}

// cancelled tells whether the player went back this frame: Escape (the 3DS's B).
func cancelled(keys *engine.Input) bool { return keys.Pressed(engine.KeyEscape) }

// stage gets a scene ready for one of the game's screens: the arena's space built in it, the
// camera seeing as far as the sky.
func stage(context *engine.Context) *arena.World {
	world := arena.NewWorldIn(context.Scene.World())
	shownBackdrop = backdrop(context)
	engine.SetSecondScreen(nil) // a screen that shows a scene there (the hangar, the radar) puts it back
	if camera, found := context.Scene.MainCamera(); found {
		if settings, isCamera := camera.Camera(); isCamera {
			settings.Near, settings.Far = arena.CameraNear, arena.CameraFar
		}
	}
	return world
}

// look puts the scene's main camera where view is, seeing as wide.
func look(context *engine.Context, view models.Camera) {
	if camera, found := context.Scene.MainCamera(); found {
		aim(camera, view)
	}
}

// Title is the first screen, as the 3DS opens: the name large over space, a ship crossing it with
// a barrel roll, a prompt that blinks.
type Title struct {
	world *arena.World
	now   float64
}

// Start builds the title's space; on the 3DS the touch screen shows the prompt.
func (title *Title) Start(context *engine.Context) {
	title.world = stage(context)
	drawInterface(nil)
	drawOver(func(drawing *draw2d.List, width, height float32) {
		paintScaled(drawing, width, height, title.paint)
	})
	if hasTouchScreen {
		showTouchScreen(func(drawing *draw2d.List) { touchscreen.TitlePrompt(drawing, textMeasurer(), title.now) })
	}
}

// Update flies the title's ship, and goes on to the modes when the player says so.
func (title *Title) Update(context *engine.Context) {
	title.now = context.Time
	darkenSky(context.Scene)
	var view models.Camera
	title.world.ShowTitle(context.Time, &view)
	look(context, view)
	if confirmed(context.Input) || touchInput(context.Input).Pressed&touchscreen.ButtonTouch != 0 {
		playSound(func(sound *soundSet) { sound.Play(sound.Choose, 1) })
		context.ChangeScene(modeScene)
	}
}

func (title *Title) paint(drawing *draw2d.List, width, height float32) {
	topHeading(drawing, width, "Ships 3D", platformCaption)
	if !hasTouchScreen {
		prompt(drawing, "PRESS ENTER OR A", width, height*0.75, title.now)
	}
}

// The modes a player chooses between; the assault isn't in this client yet.
var modes = []touchscreen.Mode{
	{Name: "FREE FOR ALL", Text: "Every pilot for themselves", Icon: touchscreen.FreeForAllIcon, Open: true},
	{Name: "ASSAULT", Text: "Coming soon", Icon: touchscreen.AssaultIcon},
}

// Modes is choosing the game's mode: a window with a big icon for each, the star on the one
// chosen; the arrows (the cross) choose, Enter (A) goes on, Escape (B) back to the title.
type Modes struct {
	world  *arena.World
	chosen int
	now    float64
}

// Start builds the space behind the window; on the 3DS the modes are on the touch screen.
func (choice *Modes) Start(context *engine.Context) {
	choice.world = stage(context)
	drawInterface(nil)
	drawOver(func(drawing *draw2d.List, width, height float32) {
		paintScaled(drawing, width, height, choice.paint)
	})
	if hasTouchScreen {
		showTouchScreen(func(drawing *draw2d.List) {
			touchscreen.ModeChoice(drawing, textMeasurer(), modes, choice.chosen, choice.now)
		})
	}
}

// Update moves the choice, and goes on or back.
func (choice *Modes) Update(context *engine.Context) {
	choice.now = context.Time
	darkenSky(context.Scene)
	var view models.Camera
	choice.world.ShowTitle(context.Time, &view)
	choice.world.HideTitle()
	look(context, view)
	keys := context.Input
	read := touchInput(keys)
	touched := choice.touched(read)
	if read.Pressed&(touchscreen.ButtonLeft|touchscreen.ButtonRight) != 0 {
		choice.chosen = (choice.chosen + 1) % len(modes)
		playSound(func(sound *soundSet) { sound.Play(sound.Move, 1) })
	}
	switch {
	case (confirmed(keys) || touched) && modes[choice.chosen].Open:
		playSound(func(sound *soundSet) { sound.Play(sound.Choose, 1) })
		context.ChangeScene(shipScene)
	case cancelled(keys):
		playSound(func(sound *soundSet) { sound.Play(sound.Back, 1) })
		context.ChangeScene(titleScene)
	}
}

// touched chooses the mode whose card a touch landed on, and tells whether one did.
func (choice *Modes) touched(read touchscreen.Input) bool {
	if read.Pressed&touchscreen.ButtonTouch == 0 {
		return false
	}
	index := touchscreen.ModeTouched(read.TouchX, read.TouchY, len(modes))
	if index < 0 {
		return false
	}
	choice.chosen = index
	return true
}

// The mode window's measures, in pixels.
const (
	modeWindowWidth  = 300
	modeWindowHeight = 150
	modeCardWidth    = 136
	modeIconPixel    = 3
)

func (choice *Modes) paint(drawing *draw2d.List, width, height float32) {
	if hasTouchScreen {
		mode := modes[choice.chosen]
		topHeading(drawing, width, "Game mode", mode.Name)
		return
	}
	area := draw2d.Rect{X: (width - modeWindowWidth) / 2, Y: (height - modeWindowHeight) / 2, Width: modeWindowWidth, Height: modeWindowHeight}
	window(drawing, area, "GAME MODE")
	for index, mode := range modes {
		card := draw2d.Rect{X: area.X + 10 + float32(index)*(modeCardWidth+8), Y: area.Y + 30, Width: modeCardWidth, Height: modeWindowHeight - 40}
		ink, iconColour := colours.soft, colours.dim
		if index == choice.chosen {
			pixelart.Stepped(drawing, card, 1, colours.chosen)
			cursor(drawing, card.X+4, card.Y+4, choice.now)
			ink, iconColour = colours.ink, colours.accent
		}
		iconWidth := float32(len(mode.Icon[0])) * modeIconPixel
		pixelart.PaintedBitmap(drawing, mode.Icon, card.X+(card.Width-iconWidth)/2, card.Y+10, modeIconPixel,
			map[byte]draw2d.Color{'#': iconColour, '+': colours.ink})
		label := card.X + (card.Width-textWidth(mode.Name, pixelart.LabelSize, draw2d.Mono))/2
		drawing.Text(mode.Name, label, card.Y+72, pixelart.LabelSize, draw2d.Mono, ink)
		text := card.X + (card.Width-textWidth(mode.Text, pixelart.TextSize, draw2d.Regular))/2
		drawing.Text(mode.Text, text, card.Y+92, pixelart.TextSize, draw2d.Regular, colours.soft)
	}
	hints(drawing, "< > CHOOSE   ENTER/A OK   ESC/B BACK", width, height)
}
