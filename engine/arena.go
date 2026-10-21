package main

import (
	"math"
	"strconv"

	"github.com/akrck02/go-engine/engine"
	"github.com/akrck02/go-renderer/draw2d"

	"github.com/jonbul/ships3D-3ds/hud"
	"github.com/jonbul/ships3D-3ds/pixelart"
	"github.com/jonbul/ships3D-3ds/radar"
	"github.com/jonbul/ships3D-3ds/touchscreen"
)

// The arena's screens: what the menus do with the keys, and what each screen shows.

// pressMenus acts on the menus' keys; it tells whether the player chose to leave the match.
func (pilot *Pilot) pressMenus(context *engine.Context) bool {
	if !hasTouchScreen {
		return pilot.pressPause(context)
	}
	keys := context.Input
	read := touchInput(keys)
	state := pilot.played.Game
	switch {
	case pausePressed(keys):
		pilot.menu.toggle()
	case pilot.menu.open:
		pilot.menu.press(pilot.pages(), read, pilot.now)
	case !state.Welcomed:
		return pilot.pressJoining(read)
	case read.Pressed&touchscreen.ButtonUp != 0:
		pilot.zoom = max(pilot.zoom-1, 0)
	case read.Pressed&touchscreen.ButtonDown != 0:
		pilot.zoom = min(pilot.zoom+1, len(radar.Ranges)-1)
	}
	return false
}

// pressJoining acts on the window shown until the server welcomes us; it tells whether the
// player went back to the ships.
func (pilot *Pilot) pressJoining(read touchscreen.Input) bool {
	read.Closes = true
	sound, closed := pilot.joining.Press([]touchscreen.Page{pilot.joiningPage()}, read, pilot.now)
	playMenuSound(sound)
	return closed
}

// pressWithoutMatch is the keys when the match couldn't start or has ended: B (Escape) goes back
// to the ships.
func (pilot *Pilot) pressWithoutMatch(context *engine.Context) {
	if hasTouchScreen && pilot.pressJoining(touchInput(context.Input)) || cancelled(context.Input) || pilot.leaving {
		pilot.leave(context)
	}
}

// pages are the game's menu: the pilot (with the way out of the game), the figures, the settings
// and the account.
func (pilot *Pilot) pages() []touchscreen.Page {
	state := pilot.played.Game
	motion := touchscreen.Motion{Roll: turning.Roll, Pitch: turning.Pitch, Yaw: turning.Yaw, TiltPitch: tilted.Pitch, TiltSteering: tilted.Steer}
	return []touchscreen.Page{
		touchscreen.PilotPage(state, touchscreen.Pilot{
			Resume:     pilot.menu.toggle,
			ChooseShip: func() { pilot.leaving = true },
			Quit:       engine.Quit,
		}),
		touchscreen.FiguresPage(pilot.figures.Shown, float64(state.Settings.TickRate), motion),
		settingsPage(),
		accountPage(func() { pilot.leaving = true }), // another server or player: join again
	}
}

// joiningPage is the window shown until the server welcomes us: how the connection goes, and the
// way back.
func (pilot *Pilot) joiningPage() touchscreen.Page {
	rows := touchscreen.StatusRows(pilot.connectionStatus(), pilot.problem())
	rows = append(rows, touchscreen.Row{}, touchscreen.Row{Label: "Back to the ships", Activate: func() { pilot.leaving = true }})
	return touchscreen.Page{Title: touchscreen.JoiningTitle, Small: true, Rows: rows}
}

// connectionStatus says how the connection to the server goes.
func (pilot *Pilot) connectionStatus() string {
	switch {
	case pilot.played == nil:
		return pilot.status
	case pilot.played.Game.Welcomed:
		return "Playing on " + player.Server
	}
	return "Joining…"
}

// problem is the server's complaint, if it made one.
func (pilot *Pilot) problem() string {
	if pilot.played == nil || pilot.played.Game.Problem == "" {
		return ""
	}
	return "Server: " + pilot.played.Game.Problem
}

// paintTouchScreen draws the touch screen: the menu, the window until we are welcomed, or the
// radar's readings over the 3D radar (or the flat map).
func (pilot *Pilot) paintTouchScreen(drawing *draw2d.List) {
	switch {
	case pilot.menu.open && pilot.played != nil:
		pilot.menu.paint(drawing, pilot.pages(), pilot.now)
	case pilot.played == nil || !pilot.played.Game.Welcomed:
		pilot.joining.Paint(drawing, []touchscreen.Page{pilot.joiningPage()}, "B SHIPS", pilot.now)
	default:
		shown := touchscreen.Radar{State: pilot.played.Game, Scanner: pilot.radar.scanner, Range: pilot.zoom, Flat: player.Settings.FlatRadar}
		touchscreen.PaintRadar(drawing, touchscreen.FullScreen, shown, pilot.now)
	}
}

// paintReadings paints over the 3D: the 3DS's life bar and the copilot, or the heading of a
// destroyed ship, laid out in the 3DS's pixels and larger on a larger screen.
func (pilot *Pilot) paintReadings(drawing *draw2d.List, width, height float32) {
	if pilot.played == nil || !pilot.played.Game.Welcomed {
		return
	}
	paintScaled(drawing, width, height, func(drawing *draw2d.List, width, height float32) {
		state := pilot.played.Game
		if !state.Alive {
			topHeading(drawing, width, "Destroyed", pilot.deathCaption())
		} else {
			hud.Paint(drawing, state, pilot.now, width, height)
		}
		if pilot.paused {
			pilot.paintPause(drawing, width, height)
		}
	})
}

// deathCaption says who destroyed us, and when we may fly again.
func (pilot *Pilot) deathCaption() string {
	state := pilot.played.Game
	caption := ""
	if killer := state.KilledBy(); killer != "" {
		caption = "by " + killer + "  ·  "
	}
	if state.CanRespawn(pilot.now) {
		return caption + "A to fly again"
	}
	return caption + "back in " + strconv.Itoa(int(math.Ceil(state.RespawnIn(pilot.now)))) + " s"
}

// pressPause is the desktop's pause menu, over the game: Escape opens it, the arrows choose,
// Enter goes; it tells whether the player chose to leave.
func (pilot *Pilot) pressPause(context *engine.Context) bool {
	keys := context.Input
	if pausePressed(keys) {
		pilot.paused, pilot.resume = !pilot.paused, true
		playSound(func(sound *soundSet) { sound.Play(sound.Choose, 1) })
		return false
	}
	if !pilot.paused {
		return false
	}
	if keys.Pressed(engine.KeyUp) || keys.Pressed(engine.KeyDown) {
		pilot.resume = !pilot.resume
		playSound(func(sound *soundSet) { sound.Play(sound.Move, 1) })
	}
	if !confirmed(keys) {
		return false
	}
	if pilot.resume {
		pilot.paused = false
		playSound(func(sound *soundSet) { sound.Play(sound.Back, 1) })
		return false
	}
	return true
}

// The pause menu's measures, in pixels.
const (
	pauseWidth  = 160
	pauseHeight = 96
	pauseRow    = 22
)

// paintPause draws the pause menu: Resume, and Leave for the ships.
func (pilot *Pilot) paintPause(drawing *draw2d.List, width, height float32) {
	area := draw2d.Rect{X: (width - pauseWidth) / 2, Y: (height - pauseHeight) / 2, Width: pauseWidth, Height: pauseHeight}
	window(drawing, area, "PAUSED")
	for index, label := range []string{"Resume", "Leave"} {
		top := area.Y + 32 + float32(index)*pauseRow
		if (index == 0) == pilot.resume {
			pixelart.Stepped(drawing, draw2d.Rect{X: area.X + 8, Y: top, Width: area.Width - 16, Height: pauseRow - 4}, 1, colours.chosen)
			cursor(drawing, area.X+14, top+5, pilot.now)
		}
		drawing.Text(label, area.X+32, top+14, pixelart.TextSize, draw2d.Regular, colours.ink)
	}
}
