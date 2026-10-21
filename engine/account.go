//go:build !nintendo3ds

package main

import ui "github.com/akrck02/go-ui"

// The account form's fields, kept between frames as the player types.
var form = struct {
	server, email, password, status string
}{}

// The form's place on the window, in pixels from its top left corner, and its width.
const (
	formLeft  = 16
	formTop   = 16
	formWidth = 300
)

// accountForm is the hangar's account window on the desktop, with a mouse and a keyboard: the
// server, and logging in or out. changed runs when either changes, for the ships to load again.
func accountForm(context *ui.Context, changed func()) {
	if form.server == "" {
		form.server, form.email = player.Server, player.Email
	}
	options := ui.PanelOptions{X: formLeft, Y: formTop, Width: ui.Fixed(formWidth), Locked: true}
	context.PanelWith("Account", options, func() {
		context.Caption("Server")
		context.TextFieldWith("server", &form.server, ui.TextFieldOptions{Placeholder: realServer, MaxLength: maximumServerLength})
		if form.server != player.Server && context.Button("Use this server") {
			form.status = errorText(useServer(form.server))
			changed()
		}
		if player.LoggedIn() {
			context.Text("Playing as "+player.Username, ui.TextStyle{})
			if context.Button("Log out") {
				logOut()
				changed()
			}
		} else {
			accountFields(context, changed)
		}
		if form.status != "" {
			context.Caption(form.status)
		}
	})
}

// accountFields are the email, the password and Log in, for a guest.
func accountFields(context *ui.Context, changed func()) {
	context.Caption("Log in to fly your own designs")
	context.TextFieldWith("email", &form.email, ui.TextFieldOptions{Placeholder: "Email", MaxLength: maximumEmailLength})
	context.TextFieldWith("password", &form.password, ui.TextFieldOptions{Placeholder: "Password", MaxLength: maximumPasswordLength, Secret: true})
	if context.Button("Log in") {
		form.status = errorText(logIn(form.email, form.password))
		form.password = ""
		changed()
	}
}

func errorText(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
