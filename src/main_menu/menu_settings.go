package main_menu

import (
	"fmt"
	"game/game_host"
	"kaiju/engine/ui/markup"
	"kaiju/engine/ui/markup/document"
	"kaiju/platform/profiler/tracing"
	"log/slog"
	"strconv"
)

const (
	settingsMenuUI = "2a55e2a9-05ff-4d20-b1a3-845c9798c4e6"
	fullScreenText = "Full Screen"
)

type menuSettingsState struct {
	doc          *document.Document
	menu         *menuStartupState
	soundPercent *document.Element
	musicPercent *document.Element
	SFXVolume    float32
	MusicVolume  float32
	Resolution   string
	Resolutions  []string
}

func (s *menuSettingsState) initialize(menu *menuStartupState) {
	defer tracing.NewRegion("menuSettingsState.initialize").End()
	s.menu = menu
	host := s.menu.host.Value()
	g := host.Game().(*game_host.GameHost)
	a := host.Audio()
	s.SFXVolume = a.SoundVolume()
	s.MusicVolume = a.MusicVolume()
	s.Resolutions = []string{
		fullScreenText,
		"1280 x 720",
		"944 x 500",
	}
	hw := host.Window
	if hw.IsFullScreen() {
		s.Resolution = fullScreenText
	} else {
		s.Resolution = fmt.Sprintf("%d x %d", hw.Width(), hw.Height())
	}
	var err error
	s.doc, err = markup.DocumentFromHTMLAsset(&g.UiMan, settingsMenuUI, s,
		map[string]func(*document.Element){
			"changeSFXVolume":     s.changeSFXVolume,
			"submitSFXVolume":     s.submitSFXVolume,
			"submitMusicVolume":   s.submitMusicVolume,
			"clickBack":           s.clickBack,
			"setWindowResolution": s.setWindowResolution,
		})
	if err != nil {
		slog.Error("failed to load the main menu UI", "error", err)
		return
	}
	s.soundPercent, _ = s.doc.GetElementById("soundPercent")
	s.musicPercent, _ = s.doc.GetElementById("musicPercent")
	s.hide()
}

func (s *menuSettingsState) show() {
	s.doc.Activate()
}

func (s *menuSettingsState) hide() {
	s.doc.Deactivate()
}

func (s *menuSettingsState) destroy() {
	s.doc.Destroy()
	s.menu = nil
}

func (s *menuSettingsState) changeSFXVolume(e *document.Element) {
	defer tracing.NewRegion("menuSettingsState.changeSFXVolume").End()
	v := e.UI.ToSlider().Value()
	s.soundPercent.InnerLabel().SetText(strconv.Itoa(int(v*100)) + "%")
}

func (s *menuSettingsState) submitSFXVolume(e *document.Element) {
	defer tracing.NewRegion("menuSettingsState.submitSFXVolume").End()
	host := s.menu.host.Value()
	if host != nil {
		host.Audio().SetSoundVolume(e.UI.ToSlider().Value())
		playButtonSound(s.menu.host.Value(), s.menu.btnSound)
	}
}

func (s *menuSettingsState) submitMusicVolume(e *document.Element) {
	defer tracing.NewRegion("menuSettingsState.submitMusicVolume").End()
	v := e.UI.ToSlider().Value()
	s.musicPercent.InnerLabel().SetText(strconv.Itoa(int(v*100)) + "%")
	host := s.menu.host.Value()
	if host != nil {
		host.Audio().SetMusicVolume(e.UI.ToSlider().Value())
	}
}

func (s *menuSettingsState) clickBack(e *document.Element) {
	defer tracing.NewRegion("menuSettingsState.clickBack").End()
	playButtonSound(s.menu.host.Value(), s.menu.btnSound)
	s.hide()
	s.menu.show()
}

func (s *menuSettingsState) setWindowResolution(e *document.Element) {
	defer tracing.NewRegion("menuSettingsState.setWindowResolution").End()
	host := s.menu.host.Value()
	if host == nil {
		return
	}
	str := e.UI.ToSelect().Value()
	if str == fullScreenText {
		host.RunOnMainThread(host.Window.SetFullscreen)
	} else {
		var w, h int
		fmt.Sscanf(str, "%d x %d", &w, &h)
		host.RunOnMainThread(func() {
			host.Window.SetWindowed(w, h)
		})
	}
}
