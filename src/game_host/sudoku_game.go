package game_host

import (
	"fmt"
	"kaiju/engine"
	"kaiju/engine/collision"
	"kaiju/engine/systems/console"
	"kaiju/matrix"
	"log/slog"
)

const (
	cellDistance = 0.91
	cellExtent   = 0.353757
)

type SudokuGame struct {
	host     *engine.Host
	CellGrid [81]collision.AABB
	updateId engine.UpdateId
}

func (g *SudokuGame) Initialize(host *engine.Host) {
	g.host = host
	for y := range 9 {
		for x := range 9 {
			pos := matrix.NewVec3(float32(x-4)*cellDistance,
				0, float32(y-4)*cellDistance)
			g.CellGrid[y*9+x] = collision.AABB{
				Center: pos,
				Extent: matrix.NewVec3(cellExtent, 0.1, cellExtent),
			}
		}
	}
}

func (g *SudokuGame) Reset() {
	// TODO:  This should reset all the squares in the cells
}

func (g *SudokuGame) Start() {
	if g.updateId.IsValid() {
		slog.Error("the game has already been started")
		return
	}
	g.updateId = g.host.Updater.AddUpdate(g.update)
}

func (g *SudokuGame) Stop() {
	g.host.Updater.RemoveUpdate(&g.updateId)
}

func (g *SudokuGame) update(float64) {
	c := g.host.Window.Cursor
	if c.Pressed() {
		r := g.host.PrimaryCamera().RayCast(c.Position())
		for i := range g.CellGrid {
			if _, ok := g.CellGrid[i].RayHit(r); ok {
				slog.Info("hit the cell", "index", i)
				console.For(g.host).Write(fmt.Sprintf("Hit the cell: %d", i))
				break
			}
		}
	}
}
