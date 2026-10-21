//go:build nintendo3ds

package main

import (
	"github.com/akrck02/go-3ds/app"
	"github.com/akrck02/go-renderer/ctr"
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"

	"github.com/jonbul/ships3D-3ds/profile"
	"github.com/jonbul/ships3D-3ds/touchscreen"
)

// figuresPageIndex is where the figures are in the game's menu, which reads the heap's size
// while they show.
const figuresPageIndex = 1

// build lays out the touch screen for what it shows now (package touchscreen paints it).
func (session *session) build(context *ui.Context) {
	layout := ui.Layout{Width: ui.Fixed(touchscreen.Width), Height: ui.Fixed(touchscreen.Height)}
	context.Paint(layout, func(area draw2d.Rect, drawing *draw2d.List) {
		switch {
		case session.screen == titleScreen && session.confirming == noChange:
			touchscreen.TitlePrompt(drawing, session.measurer, session.now)
		case session.screen == shipScreen && !session.menuOpen && session.confirming == noChange:
			touchscreen.ShipChoice(drawing, session.measurer, session.status)
		case session.screen == gameScreen && session.game.Welcomed && !session.menuOpen:
			shown := touchscreen.Radar{State: session.game, Scanner: session.scanner, Range: session.mapRange, Flat: session.profile.Settings.FlatRadar}
			touchscreen.PaintRadar(drawing, area, shown, session.now)
		default:
			pages, hint := session.pages()
			session.menu.Paint(drawing, pages, hint, session.now)
		}
	})
}

// menuInput is what the touch screen's menus read of a frame.
func menuInput(frame app.Frame) touchscreen.Input {
	return touchscreen.Input{Pressed: menuButtons(frame.Pressed), TouchX: float32(frame.Reading.TouchX), TouchY: float32(frame.Reading.TouchY)}
}

// menuButtons are the console's buttons that the menus read, in their own bits.
func menuButtons(pressed ctr.Buttons) touchscreen.Buttons {
	pairs := []struct {
		console ctr.Buttons
		menu    touchscreen.Buttons
	}{
		{ctr.ButtonUp, touchscreen.ButtonUp}, {ctr.ButtonDown, touchscreen.ButtonDown},
		{ctr.ButtonLeft, touchscreen.ButtonLeft}, {ctr.ButtonRight, touchscreen.ButtonRight},
		{ctr.ButtonA, touchscreen.ButtonA}, {ctr.ButtonB, touchscreen.ButtonB},
		{ctr.ButtonL, touchscreen.ButtonL}, {ctr.ButtonR, touchscreen.ButtonR},
		{ctr.ButtonTouch, touchscreen.ButtonTouch},
	}
	var buttons touchscreen.Buttons
	for _, pair := range pairs {
		if pressed&pair.console != 0 {
			buttons |= pair.menu
		}
	}
	return buttons
}

// pages are the menu's pages shown now, and the hint for its buttons.
func (session *session) pages() ([]touchscreen.Page, string) {
	switch {
	case session.confirming != noChange:
		return []touchscreen.Page{touchscreen.RestartWarning(session.confirm, func() { session.confirming = noChange })}, touchscreen.ChoiceHint
	case session.screen == enteringScreen:
		page := session.enteringPage()
		if session.asking == askEnter {
			return []touchscreen.Page{page}, touchscreen.WaitingHint
		}
		return []touchscreen.Page{page}, "A OK"
	case session.screen == shipScreen:
		return []touchscreen.Page{session.accountPage(), session.settingsPage(), touchscreen.ControlsPage(touchscreen.ControlLines)}, touchscreen.MenuHint
	case session.menuOpen:
		return []touchscreen.Page{session.pilotPage(), session.figuresPage(), session.settingsPage(), session.accountPage()}, touchscreen.MenuHint
	}
	return []touchscreen.Page{session.joiningPage()}, "B SHIPS"
}

// pressPages acts on the menu shown now, with its sounds; B closes it, or answers no.
func (session *session) pressPages(input touchscreen.Input) {
	pages, _ := session.pages()
	input.Closes = session.confirming != noChange || session.menuOpen || session.screen == gameScreen
	input.TurnsPages = session.menuOpen && session.confirming == noChange
	sound, closed := session.menu.Press(pages, input, session.now)
	session.playMenuSound(sound)
	if !closed {
		return
	}
	switch {
	case session.confirming != noChange:
		session.confirming = noChange
	case session.menuOpen:
		session.menuOpen = false
	case session.screen == gameScreen:
		session.leave("")
	}
}

// playMenuSound plays the sound a menu asked for.
func (session *session) playMenuSound(sound touchscreen.Sound) {
	played := session.sounds
	switch sound {
	case touchscreen.Moved:
		played.Play(played.Move, 1)
	case touchscreen.Chose:
		played.Play(played.Choose, 1)
	case touchscreen.Went:
		played.Play(played.Back, 1)
	}
}

// statusRows are the connection's status and the server's complaint, wrapped to the window.
func (session *session) statusRows() []touchscreen.Row {
	return touchscreen.StatusRows(session.status, session.problem())
}

func (session *session) problem() string {
	if session.game.Problem == "" {
		return ""
	}
	return "Server: " + session.game.Problem
}

// enteringPage says the server is being reached, or why it couldn't be, with a way out.
func (session *session) enteringPage() touchscreen.Page {
	rows := session.statusRows()
	if session.asking != askEnter {
		rows = append(rows, touchscreen.Row{},
			touchscreen.Row{Label: "Try again", Activate: session.enter},
			touchscreen.Row{Label: "Server and account", Activate: func() {
				session.menuOpen, session.screen = true, shipScreen
				session.menu.Reset() // the menu opens new, on its first page: the account
			}})
	}
	return touchscreen.Page{Title: touchscreen.ConnectingTitle, Small: true, Rows: rows}
}

// joiningPage shows while the game's connection is made, before the server's welcome.
func (session *session) joiningPage() touchscreen.Page {
	rows := append(session.statusRows(), touchscreen.Row{},
		touchscreen.Row{Label: "Back to the ships", Activate: func() { session.leave("") }})
	return touchscreen.Page{Title: touchscreen.JoiningTitle, Small: true, Rows: rows}
}

// pilotPage is the game's own page: life and score, the ways out, and the others' scores.
func (session *session) pilotPage() touchscreen.Page {
	return touchscreen.PilotPage(session.game, touchscreen.Pilot{
		Resume:     func() { session.menuOpen = false },
		ChooseShip: func() { session.leave("") },
		Quit:       func() { session.quitting = true },
	})
}

// settingsPage offers how the ship flies, the radar shows and the game sounds; each change is
// remembered once the menu closes.
func (session *session) settingsPage() touchscreen.Page {
	change := func(apply func(settings *profile.Settings)) {
		changed := session.profile.Settings
		apply(&changed)
		session.rememberLater(func(remembered *profile.Profile) { remembered.Settings = changed })
		session.applyVolumes()
	}
	return touchscreen.SettingsPage(session.profile.Settings, change, session.tilt.Centre)
}

// accountPage shows the server and the account, each with its change, which restarts the game.
func (session *session) accountPage() touchscreen.Page {
	return touchscreen.AccountPage(touchscreen.Account{
		Server:       session.profile.Server,
		Username:     session.accountName(),
		ChangeServer: func() { session.ask(changeServer) },
		LogIn:        func() { session.ask(changeLogin) },
		LogOut:       func() { session.ask(changeLogout) },
		Status:       session.statusRows(),
	})
}

// accountName is the logged in player's name, or nothing for a guest.
func (session *session) accountName() string {
	if session.account.Token == "" {
		return ""
	}
	return session.account.Username
}

// figuresPage shows how the game runs: the frame, the drawing, the network and the tilt.
func (session *session) figuresPage() touchscreen.Page {
	gyroscope := session.gyroscope
	motion := touchscreen.Motion{Roll: gyroscope.Roll, Pitch: gyroscope.Pitch, Yaw: gyroscope.Yaw, TiltPitch: session.tilt.Pitch, TiltSteering: session.tilt.Steer}
	return touchscreen.FiguresPage(session.figures.Shown, float64(session.game.Settings.TickRate), motion)
}
