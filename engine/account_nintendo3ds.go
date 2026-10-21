//go:build nintendo3ds

package main

import (
	"github.com/akrck02/go-3ds/keyboard"
	"github.com/akrck02/go-engine/engine"
	ui "github.com/akrck02/go-ui"
)

// status is what the account screen says after a change: why a login failed.
var status string

// accountForm is the hangar's account on the 3DS's bottom screen: the server and the player, with
// buttons that open the system keyboard. The keyboard takes both screens, so it opens between
// frames (engine.BetweenFrames). changed runs when either changes, for the ships to load again.
func accountForm(context *ui.Context, changed func()) {
	context.Caption("Server")
	context.Text(player.Server, ui.TextStyle{})
	if context.Button("Change server") {
		engine.BetweenFrames(func() {
			typed, confirmed := keyboard.Ask(keyboard.Text, "Server: address:port, or https://name:port", player.Server, maximumServerLength)
			if confirmed && typed != player.Server {
				status = errorText(useServer(typed))
				changed()
			}
		})
	}
	context.Caption("Account")
	if player.LoggedIn() {
		context.Text(player.Username, ui.TextStyle{})
		if context.Button("Log out") {
			engine.BetweenFrames(func() {
				logOut()
				changed()
			})
		}
	} else if context.Button("Log in") {
		engine.BetweenFrames(func() { askLogin(changed) })
	}
	if status != "" {
		context.Caption(status)
	}
}

// askLogin asks for the email and the password with the system keyboard, and logs in.
func askLogin(changed func()) {
	email, confirmed := keyboard.Ask(keyboard.Text, "Email", player.Email, maximumEmailLength)
	if !confirmed {
		return
	}
	password, confirmed := keyboard.Ask(keyboard.Password, "Password", "", maximumPasswordLength)
	if !confirmed {
		return
	}
	status = errorText(logIn(email, password))
	changed()
}

func errorText(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
