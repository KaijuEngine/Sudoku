package game_host

import (
	"kaiju/engine"
	"kaiju/engine/ui"
)

type GameHost struct {
	UiMan ui.Manager
	Game  SudokuGame
}

func NewGameHost(host *engine.Host) *GameHost {
	g := &GameHost{}
	g.UiMan.Init(host)
	g.Game.Initialize(host)
	return g
}
