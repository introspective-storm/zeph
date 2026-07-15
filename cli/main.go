package main

import (
	"cli/config"
	"cli/picker"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2" // Modern v2 import
)

type rootModel struct {
	state        AppState
	menu         listModel
	configPicker picker.Model
	dataPicker   picker.Model
	activeConfig *config.ProjectConfig
	filepicker   picker.Model
	selected     string
	showPicker   bool
	pickerError  string
}

func initialModel() rootModel {
	const pickerHeight = 15
	return rootModel{
		state:        StateMenu,
		menu:         mainMenu(), //the
		configPicker: picker.ConfigPicker(pickerHeight),
		dataPicker:   picker.DataPicker(pickerHeight),
		activeConfig: &config.ProjectConfig{},
	}
}

func (m rootModel) Init() tea.Cmd {
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.menu = m.menu.SetSize(msg.Width, msg.Height)
	}

	switch m.state {
	case StateMenu:
		m.menu, cmd = m.menu.Update(msg)
	}
	return m, cmd
}

func (m rootModel) View() tea.View {
	var content string

	switch m.state {
	case StateMenu:
		content = m.menu.View()
	}
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting Tahoe: %v\n", err)
		os.Exit(1)
	}
}
