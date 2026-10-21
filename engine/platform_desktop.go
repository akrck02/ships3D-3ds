//go:build !nintendo3ds

package main

// hasTouchScreen is false on the desktop: one screen shows the game, its prompts and its menus.
//
// TODO(akrck02): show the touch screen in a window over the game, as the 3DS's bottom screen
// (proposal 3, T6).
const hasTouchScreen = false

// platformCaption is the title's caption, under the game's name.
const platformCaption = "a space dogfight"

// paintsSky is true on the desktop: its renderer paints a day's sky wherever nothing covers the
// far plane, with the sun's glow, so the game darkens it (darkenSky) and the backdrop covers it.
const paintsSky = true
