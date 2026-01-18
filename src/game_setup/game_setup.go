package game_setup

import (
	"game/game_host"
	"kaiju/engine"
)

func init() {
	engine.RegisterEntityData("game_setup.GameSetup", GameSetup{})
}

type GameSetup struct {
	// Setup your POD
}

func (d GameSetup) Init(_ *engine.Entity, host *engine.Host) {
	g := host.Game().(*game_host.GameHost)
	g.Game.Reset()
	g.Game.Start()
}
