package main

import (
	"cli/config"
	"cli/picker"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2" // Modern v2 import[cite: 7]
)

// 1. Define your AppState enum

type rootModel struct {
	state        AppState
	menu         listModel
	wizard       config.Wizard
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
		menu:         mainMenu(),
		wizard:       config.NewWizard(pickerHeight),
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
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// If we are in the main menu, 'q' quits
		if m.state == StateMenu && msg.String() == "q" {
			return m, tea.Quit //[cite: 7]
		}
		if m.state == StateWizard && msg.String() == "q" {
			m.state = StateMenu
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.menu = m.menu.SetSize(msg.Width, msg.Height) //[cite: 7]
	}

	// 2. Route input to the active screen state
	switch m.state {
	case StateMenu:
		m.menu, cmd = m.menu.Update(msg) //[cite: 7]

		// Check if the user hit Enter on a menu item
		if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "enter" {
			if selectedItem, ok := m.menu.SelectedItem().(item); ok { //[cite: 8]
				if selectedItem.GetAction() == "wizard" { //[cite: 8]
					// Start a clean wizard and switch states!
					m.wizard = config.NewWizard(15)
					m.state = StateWizard
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
	}

	return m, cmd
}

func (m rootModel) View() tea.View {
	var content string

	// 3. Render the correct screen based on active state
	switch m.state {
	case StateMenu:
		content = m.menu.View() //[cite: 7]
	case StateWizard:
		content = m.wizard.View()
	}

	v := tea.NewView(content) //[cite: 7]
	v.AltScreen = true        //[cite: 7]
	return v                  //[cite: 7]
}

func main() {
	p := tea.NewProgram(initialModel()) //[cite: 7]
	if _, err := p.Run(); err != nil {  //[cite: 7]
		fmt.Printf("Error starting Tahoe: %v\n", err) //[cite: 7]
		os.Exit(1)                                    //[cite: 7]
	}
}
