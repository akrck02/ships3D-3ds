package main

import (
	"github.com/akrck02/go-renderer/draw2d"

	"github.com/jonbul/ships3D-3ds/touchscreen"
)

// touchMenu is the menu that START opens on the touch screen, over the ship screen or the game,
// as in the hand-written program: pages with tabs that L and R (or a touch) go through.
type touchMenu struct {
	open  bool
	shown touchscreen.Menu
}

// toggle opens or closes the menu, with its sound; closing it saves the settings.
func (menu *touchMenu) toggle() {
	menu.open = !menu.open
	if menu.open {
		menu.shown.Reset() // it opens new, on its first page
		playSound(func(sound *soundSet) { sound.Play(sound.Choose, 1) })
		return
	}
	saveSettings()
	playSound(func(sound *soundSet) { sound.Play(sound.Back, 1) })
}

// press acts on the menu's pages; B closes it.
func (menu *touchMenu) press(pages []touchscreen.Page, read touchscreen.Input, now float64) {
	read.Closes, read.TurnsPages = true, true
	sound, closed := menu.shown.Press(pages, read, now)
	if closed {
		menu.toggle()
		return
	}
	playMenuSound(sound)
}

// paint draws the menu's page.
func (menu *touchMenu) paint(drawing *draw2d.List, pages []touchscreen.Page, now float64) {
	menu.shown.Paint(drawing, pages, touchscreen.MenuHint, now)
}
