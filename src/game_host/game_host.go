package game_host

import (
	"kaijuengine.com/engine"
	"kaijuengine.com/engine/stages"
	"kaijuengine.com/engine/ui"
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

func (g *GameHost) MainLoaded(host *engine.Host, loadResult stages.LoadResult) {}
