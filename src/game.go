package main

import (
	"encoding/json"
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
	s := stages.Stage{}
	if build.Debug && !klib.IsMobile() {
		j := stages.StageJson{}
		if err := json.Unmarshal(stageData, &j); err != nil {
			slog.Error("failed to decode the entry point stage 'main'", "error", err)
			host.Close()
			return
		}
		s.FromMinimized(j)
	} else {
		if s, err = stages.ArchiveDeserializer(stageData); err != nil {
			slog.Error("failed to deserialize the entry point stage", "stage", startStage, "error", err)
			host.Close()
			return
		}
	}
	host.SetGame(game_host.NewGameHost(host))
	s.Load(host)
}

func getGame() bootstrap.GameInterface { return Game{} }
