package main

import (
	"game/game_host"
	"kaiju/bootstrap"
	"kaiju/build"
	"kaiju/engine"
	"kaiju/engine/assets"
	"kaiju/klib"
	"kaiju/engine/stages"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
)

type Game struct{}

func (Game) PluginRegistry() []reflect.Type { return []reflect.Type{} }

func (Game) ContentDatabase() (assets.Database, error) {
	if klib.IsMobile() {
		return assets.NewArchiveDatabase("game.dat", []byte(build.ArchiveEncryptionKey))
	}
	if build.Debug {
		return assets.DebugContentDatabase{}, nil
	}
	p, err := os.Executable()
	if err == nil {
		pDir := filepath.Dir(p)
		dat := filepath.Join(pDir, "game.dat")
		if _, err := os.Stat(dat); err != nil {
			if _, err := os.Stat(filepath.Join(pDir, "../kaiju")); err == nil {
				if _, err := os.Stat(filepath.Join(pDir, "../src/go.mod")); err == nil {
					if _, err := os.Stat(filepath.Join(pDir, "../build/game.dat")); err == nil {
						dat = filepath.Join(pDir, "../build/game.dat")
					}
				}
			}
		}
		return assets.NewArchiveDatabase(dat, []byte(build.ArchiveEncryptionKey))
	} else {
		return assets.NewArchiveDatabase("game.dat", []byte(build.ArchiveEncryptionKey))
	}
}

func (Game) Launch(host *engine.Host) {
	startStage := engine.LaunchParams.StartStage
	if startStage == "" {
		var err error
		startStage, err = host.AssetDatabase().ReadText(stages.EntryPointAssetKey)
		if err != nil {
			slog.Error("failed to read the entry point stage id from the asset database", "key", stages.EntryPointAssetKey, "error", err)
			host.Close()
			return
		}
	}
	stageData, err := host.AssetDatabase().Read(startStage)
	if err != nil {
		slog.Error("failed to read the entry point stage", "stage", startStage, "error", err)
		host.Close()
		return
	}
	s, err := stages.Deserialize(stageData)
	if err != nil {
		slog.Error("failed to deserialize the entry point stage", "stage", startStage, "error", err)
		host.Close()
		return
	}
	gh := game_host.NewGameHost(host)
	host.SetGame(gh)
	loadResult := s.Load(host)
	gh.MainLoaded(host, loadResult)
}

func getGame() bootstrap.GameInterface { return Game{} }
