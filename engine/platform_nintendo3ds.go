//go:build nintendo3ds

package main

// hasTouchScreen is true on the 3DS: the touch screen shows the menus and the radar, and the top
// screen only the game and its headings.
const hasTouchScreen = true

// platformCaption is the title's caption, under the game's name.
const platformCaption = "for Nintendo 3DS"

// paintsSky is false on the 3DS: it clears the screen to the sky's colour and paints nothing else.
const paintsSky = false
