package main_menu

import (
	"game/game_host"
	"kaiju/engine"
	"kaiju/engine/ui/markup"
	"kaiju/engine/ui/markup/document"
	"kaiju/platform/audio"
	"kaiju/platform/profiler/tracing"
	"log/slog"
	"weak"
)

const mainMenuUI = "f3fc7333-3038-449d-a4ce-15535c8047fd"
const buttonSoundId = "f7625352-2c7a-43d2-ac1a-a6a7d0b18f35"

func init() {
	engine.RegisterEntityData(MenuStartup{})
}

type MenuStartup struct {
	// Setup your POD
}

type menuStartupState struct {
	host     weak.Pointer[engine.Host]
	doc      *document.Document
	btnSound *audio.AudioClip
	settings menuSettingsState
}

func (d MenuStartup) Init(entity *engine.Entity, host *engine.Host) {
	defer tracing.NewRegion("MenuStartup.Init").End()
	g := host.Game().(*game_host.GameHost)
	state := &menuStartupState{
		host: weak.Make(host),
	}
	var err error
	state.doc, err = markup.DocumentFromHTMLAsset(&g.UiMan, mainMenuUI, nil,
		map[string]func(*document.Element){
			"clickPlay":     state.clickPlay,
			"clickSettings": state.clickSettings,
			"clickExit":     state.clickExit,
		})
	if err != nil {
		slog.Error("failed to load the main menu UI", "error", err)
		return
	}
	entity.OnDestroy.Add(func() {
		state.doc.Destroy()
		state.settings.destroy()
	})
	state.btnSound, err = host.Audio().LoadSound(host.AssetDatabase(), buttonSoundId)
	if err != nil {
		slog.Error("failed to load the audio clip for the button",
			"id", buttonSoundId, "error", err)
	}
	state.settings.initialize(state)
}

func (s *menuStartupState) clickPlay(e *document.Element) {
	defer tracing.NewRegion("menuStartupState.clickPlay").End()
	playButtonSound(s.host.Value(), s.btnSound)
	slog.Info("clicked on play button")
}

func (s *menuStartupState) clickSettings(e *document.Element) {
	defer tracing.NewRegion("menuStartupState.clickSettings").End()
	playButtonSound(s.host.Value(), s.btnSound)
	s.hide()
	s.settings.show()
	slog.Info("clicked on settings button")
}

func (s *menuStartupState) clickExit(e *document.Element) {
	defer tracing.NewRegion("menuStartupState.clickExit").End()
	playButtonSound(s.host.Value(), s.btnSound)
	if h := s.host.Value(); h != nil {
		h.Close()
	}
}

func (s *menuStartupState) show() {
	s.doc.Activate()
}

func (s *menuStartupState) hide() {
	s.doc.Deactivate()
}

func playButtonSound(host *engine.Host, btnSound *audio.AudioClip) {
	if host != nil && btnSound != nil {
		host.Audio().Play(btnSound)
	}
}
