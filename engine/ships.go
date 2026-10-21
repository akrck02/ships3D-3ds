package main

import (
	"math"
	"strconv"

	"github.com/akrck02/go-engine/engine"
	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/scene"

	"github.com/jonbul/ships3D-3ds/account"
	"github.com/jonbul/ships3D-3ds/arena"
	"github.com/jonbul/ships3D-3ds/flight"
	"github.com/jonbul/ships3D-3ds/pixelart"
	"github.com/jonbul/ships3D-3ds/protocol"
	"github.com/jonbul/ships3D-3ds/ship"
	"github.com/jonbul/ships3D-3ds/transport"
)

// The hangar's turning ship: how big it shows, and how fast it turns, in turns a second.
const (
	hangarShipRadius    = 7
	hangarTurnsASecond  = 0.12
	shipNameCharacters  = 18
	arrowsFromTheMiddle = 150
)

// Ships is choosing the ship in the hangar, as the browser and the 3DS show it: the ship turning
// on a platform floating in space, its name at the top, the arrows (the cross) to go through the
// player's ships, Enter (A) to launch, Escape (B) back to the modes.
type Ships struct {
	Server string `json:"server"` // address:port, or https://name:port; empty is the real server

	world   *arena.World
	hangar  *arena.Hangar
	turning *scene.Node
	ships   []protocol.Ship
	index   int
	status  string
	now     float64
}

// Start builds the hangar and asks the server for its ships.
func (choice *Ships) Start(context *engine.Context) {
	choice.world = stage(context)
	choice.hangar = arena.NewHangar()
	world := context.Scene.World()
	world.Nodes = append(world.Nodes, choice.hangar.Nodes...)
	choice.turning = &scene.Node{Name: "chosen ship", Moving: true}
	world.Nodes = append(world.Nodes, choice.turning)
	choice.ships, choice.status = fetchShips(orDefault(choice.Server, realServer))
	for index, offered := range choice.ships {
		if offered.Id == chosen.shipID {
			choice.index = index
		}
	}
	drawOver(func(drawing *draw2d.List, width, height float32) {
		paintScaled(drawing, width, height, choice.paint)
	})
}

// fetchShips is the server's ships for a guest, the built-in ones; when the server can't be
// reached, the Arrow, and why.
func fetchShips(server string) ([]protocol.Ship, string) {
	fallback := []protocol.Ship{{Id: defaultShip, Name: "Arrow"}}
	address, err := account.ParseAddress(server)
	if err != nil {
		return fallback, err.Error()
	}
	api := account.Server{Host: address.HostPort(), Dial: func() (account.Connection, error) {
		return transport.Dial(address.Host, address.Port, address.Secure)
	}}
	offered, err := api.Ships(account.Session{})
	if err != nil || len(offered.Defaults)+len(offered.Own) == 0 {
		return fallback, "The server's ships didn't load"
	}
	return append(offered.Own, offered.Defaults...), ""
}

// Update turns the ship, and goes through the ships, on or back.
func (choice *Ships) Update(context *engine.Context) {
	choice.now = context.Time
	darkenSky(context.Scene)
	choice.hangar.Show(context.Time)
	choice.world.HideTitle()
	look(context, arena.HangarCamera)
	choice.showShip(context.Time)
	keys := context.Input
	switch {
	case keys.Pressed(engine.KeyLeft):
		choice.turn(-1)
	case keys.Pressed(engine.KeyRight):
		choice.turn(1)
	case confirmed(keys):
		chosen.shipID = choice.ships[choice.index].Id
		playSound(func(sound *soundSet) { sound.Play(sound.Choose, 1) })
		context.ChangeScene(arenaScene)
	case cancelled(keys):
		playSound(func(sound *soundSet) { sound.Play(sound.Back, 1) })
		context.ChangeScene(modeScene)
	}
}

func (choice *Ships) turn(by int) {
	choice.index = (choice.index + by + len(choice.ships)) % len(choice.ships)
	playSound(func(sound *soundSet) { sound.Play(sound.Move, 1) })
}

// showShip puts the chosen ship on the platform, turning.
func (choice *Ships) showShip(now float64) {
	offered := choice.ships[choice.index]
	radius := ship.Radius(offered.Layers)
	choice.turning.Hidden = radius <= 0
	if choice.turning.Hidden {
		return
	}
	choice.turning.Mesh = choice.world.ShipMesh(offered.Id, offered.Layers, false)
	turn := flight.FromEulerXYZ(0, now*hangarTurnsASecond*2*math.Pi, 0)
	choice.turning.Transform = arena.TransformOf(arena.HangarShipPlace, turn, hangarShipRadius/radius)
}

func (choice *Ships) paint(drawing *draw2d.List, width, height float32) {
	offered := choice.ships[choice.index]
	name := []rune(offered.Name)
	if len(name) > shipNameCharacters {
		name = append(name[:shipNameCharacters-1], '…')
	}
	centred(drawing, string(name), width, 30, headingSize*1.5, draw2d.Bold, colours.ink)
	count := strconv.Itoa(choice.index+1) + " / " + strconv.Itoa(len(choice.ships))
	centred(drawing, count, width, 46, pixelart.LabelSize, draw2d.Mono, colours.soft)
	middle := height / 2
	drawing.Text("<", width/2-arrowsFromTheMiddle, middle, headingSize, draw2d.Bold, colours.accent)
	drawing.Text(">", width/2+arrowsFromTheMiddle-8, middle, headingSize, draw2d.Bold, colours.accent)
	if choice.status != "" {
		centred(drawing, choice.status, width, height-30, pixelart.TextSize, draw2d.Regular, colours.soft)
	}
	prompt(drawing, "ENTER/A LAUNCH", width, height-20, choice.now)
	hints(drawing, "< > SHIP   ESC/B BACK", width, height)
}
