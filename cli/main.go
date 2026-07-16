package main

import (
	"cli/config"
	"cli/picker"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type rootModel struct {
	state        AppState
	menu         listModel
	wizard       config.Wizard
	configPicker picker.ConfigLoader
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
		menu:         mainMenu(),
		wizard:       config.NewWizard(pickerHeight),
		configPicker: picker.NewConfigLoader(pickerHeight),
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
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// If we are in the main menu, 'q' quits
		if m.state == StateMenu && msg.String() == "q" {
			return m, tea.Quit
		}
		if m.state == StateWizard && msg.String() == "q" {
			m.state = StateMenu
			return m, nil
		}
		if m.state == StateConfigPicker && msg.String() == "q" {
			m.state = StateMenu
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.menu = m.menu.SetSize(msg.Width, msg.Height)
	}

	switch m.state {
	case StateMenu:
		m.menu, cmd = m.menu.Update(msg)

		if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "enter" {
			if selectedItem, ok := m.menu.SelectedItem().(item); ok {
				if selectedItem.GetAction() == "wizard" {
					m.wizard = config.NewWizard(15)
					m.state = StateWizard
				}
				if selectedItem.GetAction() == "load" {
					m.configPicker = picker.NewConfigLoader(15)
					m.state = StateConfigPicker
				}
			}
		}

	case StateWizard:
		m.wizard, cmd = m.wizard.Update(msg)

		// Check if the wizard successfully finished writing the config file
		if m.wizard.Done {
			// Back to the menu, or you can route straight into loading the newly created project!
			m.state = StateMenu
		}
	case StateConfigPicker:
		m.configPicker, cmd = m.configPicker.Update(msg)

		if m.configPicker.Done {
			m.selected = m.configPicker.Path
			m.state = StateMenu
		}
		return m, cmd
	}

	return m, cmd
}

func (m rootModel) View() tea.View {
	var content string

	// 3. Render the correct screen based on active state
	switch m.state {
	case StateMenu:
		content = m.menu.View()
	case StateWizard:
		content = m.wizard.View()
	case StateConfigPicker:
		content = m.configPicker.View()
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
