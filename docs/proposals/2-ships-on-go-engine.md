# 2. Ships 3D as a go-engine game

> Status: proposed · Author's approval: proposal ☐ plan ☐

## Why

Ships 3D's 3DS client (this repository) is a program written by hand:
its own loop, its own scene, its own menus. go-engine now does much of that for any game: scenes
edited in an editor, behaviours in Go, cameras, lights, physics, the bottom screen, and builds for
the desktop, Android and the 3DS from one project. Moving Ships 3D onto it makes the game easier
to change (the hangar, the arena and the HUD become scenes and behaviours), puts it on the desktop
and on phones for free, and is the first online, 3D, two-screen game to test go-engine with.

## What

| Id | Requirement | Check |
|---|---|---|
| R1 | A go-engine project, Ships 3D, opens in the editor: scenes for the title, the hangar and the arena, laid out there | The project opens with `go-engine <folder>`; the three scenes show in the editor |
| R2 | The same game runs on the desktop, the 3DS and Android, built from the editor's Build | A build for each; desktop and 3DS played against a server |
| R3 | It plays online with the browser's players, as `client3ds` does: log in, the player's ships, join, fly, fire, hits, deaths, respawns | Played against a local server with bots, and against the real one from the 3DS |
| R4 | The ships are built from their designs (the designer's parts), as the browser shows them | Test: a design's mesh has the parts' triangles; a ship seen in the hangar and the arena |
| R5 | The same experience as the browser and the 3DS today: the title with a ship crossing it, the mode, the ship in the hangar, the game, the pause menu | Played on the 3DS and the desktop, side by side with the browser |
| R6 | The arena: the chase camera, the stars and planets, the world's edge as a force field, the copilot turning the ship back, the HUD (the 3DS's life bar) | Played; the edge's behaviour checked by the existing game tests |
| R7 | On the 3DS: the 3D radar and the menu on the touch screen, stereoscopic 3D, the gyroscope | Played on a New 3DS |
| R8 | Sound: chiptune effects, music and the engine's hum, on every platform | Heard on the desktop and the 3DS |
| R9 | The game's logic is not copied: the engine game uses this repository's plain-Go packages (protocol, game, flight, ship, radar, chase, impact, chiptune, account, profile) | It lives in this repository and imports them; no copies |
| R10 | 60 frames a second on a New 3DS with a dozen players | STATS on the console |

## Out of scope

- The assault mode, ship textures, the `crash` message and the language setting (each can follow,
  in this project or in `client3ds`).
- Removing or replacing the hand-written 3DS program (`cmd/ships3ds`): it stays as it is until the
  engine game plays as well.
- An old (non-New) 3DS.
- go-engine's own docs and versions: the go-engine session looks after those.

## Design

**Where.** In this repository, as a go-engine project in `engine/` (made with `go-engine new`):
its scenes, behaviours and `main.go`, in the same module, `github.com/jonbul/ships3D-3ds`.

**What it reuses.** The repository's plain-Go packages hold the game: the protocol, the game's
state, the flight model, ships' meshes, the radar, the camera, crash effects, sounds made in code,
login and the saved profile. Being in the same module, `engine/` imports them directly.
Alternatives: copy the packages (breaks R9), or a repository of its own (another module to keep in
step).

**What the game adds.** Behaviours (`Hangar`, `Pilot`, `Arena`, `Radar`, `Copilot`, `Engine`),
the scenes, and a network piece: on the 3DS it dials with go-3ds `network` (TCP and TLS); on the
desktop and Android, with Go's `net` and `crypto/tls`. Both satisfy the transport `websocket` and
`web` already take.

**What go-engine needs**, asked of the go-engine session, which makes them in the engine (the
author's choice); the game uses the versions it releases:

| Gap | Why Ships needs it |
|---|---|
| A node whose mesh is made by the game's code (`Node.SetMesh`) | Ships are built from their designs at run time, not imported `.glb` files |
| Nodes marked as moving (go-renderer's `Node.Moving`) | Ships turn every frame; baked lighting stalls the 3DS |
| A 3D scene on the bottom screen as well as a 2D one | The 3D radar and the hangar's turning ship |
| A 2D layer over the 3D top screen, on both eyes | The life bar and the copilot |
| Sound through go-sound on every platform | R8 |
| The 3DS's gyroscope, C-stick and START in `engine.Input` | Controls the 3DS client has today |

If a gap can't be closed in go-engine soon, the game draws that piece itself with go-renderer, as
`client3ds` does, and the proposal says so.

**Alternatives not taken.** Keep extending `client3ds` alone: it works, but it stays 3DS-only and
every screen is code. Rewrite the logic inside the engine's project: duplicates what is tested.

## Decided by the author

1. It lives in jonbul's repository for the 3DS client, ships3D-3ds, in `engine/`.
2. What go-engine lacks is asked of the go-engine session.
3. The experience is the same everywhere: the browser's and the 3DS's, which already match.

Still open: whether the desktop build is for players (the login and server screens designed for a
mouse and keyboard, and the game released for the desktop as well) or only for working on the game
in the editor (players keep the browser and the 3DS).

## Tasks

To be planned once the proposal is approved.

## Result

Not started.
