//go:build !nintendo3ds

package main

// hasTouchScreen is false on the desktop: one window shows the game, and the touch screen's pages
// in the middle of it when they open (panel.go), clicked with the mouse.
const hasTouchScreen = false

// platformCaption is the title's caption, under the game's name.
const platformCaption = "a space dogfight"

// paintsSky is true on the desktop: its renderer paints a day's sky wherever nothing covers the
// far plane, with the sun's glow, so the game darkens it (darkenSky) and the backdrop covers it.
const paintsSky = true
