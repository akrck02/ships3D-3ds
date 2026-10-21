package main

import "testing"

func TestAClickLandsOnThePanelInItsOwnPixels(test *testing.T) {
	windowWidth, windowHeight = 1200, 720 // three times the top screen
	x, y, inside := onPanel(1200/2, 720/2)
	if !inside || x != 160 || y != 120 {
		test.Errorf("the window's middle lands at %v, %v (inside %v), want the panel's middle, 160, 120", x, y, inside)
	}
	if _, _, inside := onPanel(10, 10); inside {
		test.Error("a click in the window's corner lands on the panel")
	}
}
