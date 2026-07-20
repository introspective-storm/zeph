package main

import (
	"cli/config"
	"cli/picker"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type rootModel struct {
	state        AppState
	activeConfig *config.ProjectConfig
	menu         listModel
	wizard       config.Wizard
	configPicker picker.ConfigLoader
	dataPicker   picker.Model
	filepicker   picker.Model
	selected     string
	showPicker   bool
	pickerError  string
	err          error
}

func initialModel() rootModel {
	var matches []string
	const pickerHeight = 15
	entries, err := os.ReadDir(".")
	if err != nil {
		return rootModel{err: err}
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".zeph.yaml") || strings.HasSuffix(name, ".zeph.yml") {
			matches = append(matches, name)
		}
	}
	matchesLen := len(matches)
	if matchesLen > 1 {
		return rootModel{
			state:        StateClarifyMenu,
			activeConfig: &config.ProjectConfig{},
			menu:         mainMenu(),
			wizard:       config.NewWizard(pickerHeight),
			configPicker: picker.NewConfigLoaderLocked(pickerHeight),
			dataPicker:   picker.DataPicker(pickerHeight),
		}
	}
	if matchesLen == 1 {
		filePath := fmt.Sprintf("./%s", matches[0])
		project, err := config.Load(filePath)
		if err != nil {
			return rootModel{err: err}
		}
		return rootModel{
			state:        StateMenu,
			activeConfig: project,
			menu:         mainMenu().WithProject(project.Name),
			wizard:       config.NewWizard(pickerHeight),
			configPicker: picker.NewConfigLoader(pickerHeight),
			dataPicker:   picker.DataPicker(pickerHeight),
		}
	}
	return rootModel{
		state:        StateMenu,
		activeConfig: &config.ProjectConfig{},
		menu:         mainMenu(),
		wizard:       config.NewWizard(pickerHeight),
		configPicker: picker.NewConfigLoader(pickerHeight),
		dataPicker:   picker.DataPicker(pickerHeight),
	}
}

func (m rootModel) Init() tea.Cmd {
	return m.configPicker.Init()
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// if we are in the main menu or clarification menu, 'q' quits
		if m.state == StateClarifyMenu && msg.String() == "q" {
			return m, tea.Quit
		}
		if m.state == StateMenu && msg.String() == "q" {
			return m, tea.Quit
		}
		// if not, then it goes back to main menu
		if m.state == StateWizard && msg.String() == "q" {
			m.state = StateMenu
			return m, nil
		}
		if m.state == StateConfigPicker && msg.String() == "q" {
			m.state = StateMenu
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.menu = m.menu.SetSize(msg.Width, 15)
	}

	switch m.state {
	case StateClarifyMenu:
		m.configPicker, cmd = m.configPicker.Update(msg)
		if m.configPicker.Done {
			m.selected = m.configPicker.Path
			project, err := config.Load(m.selected)
			if err != nil {
				m.err = err
				m.state = StateMenu
				return m, nil
			}
			m.activeConfig = project
			m.menu = m.menu.WithProject(m.activeConfig.Name)
			m.state = StateMenu
			return m, cmd
		}
	case StateMenu:
		m.menu, cmd = m.menu.Update(msg)

		if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "enter" {
			if selectedItem, ok := m.menu.SelectedItem().(item); ok {
				if selectedItem.GetAction() == "cont" {
					// load the project state here
				}
				if selectedItem.GetAction() == "load" {
					m.configPicker = picker.NewConfigLoader(15)
					m.state = StateConfigPicker
					return m, m.configPicker.Init()
				}
				if selectedItem.GetAction() == "wizard" {
					m.wizard = config.NewWizard(15)
					m.state = StateWizard
					return m, m.wizard.Init()
				}
			}
		}
		return m, cmd

	case StateWizard:
		m.wizard, cmd = m.wizard.Update(msg)

		// Check if the wizard successfully finished writing the config file
		if m.wizard.Done {
			m.activeConfig = m.wizard.Conf
			m.state = StateMenu
			m.menu = m.menu.WithProject(m.activeConfig.Name)
		}
		return m, cmd
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
	case StateClarifyMenu:
		header := "Multiple '.zeph.yaml'/'.zeph.yml' files found, please choose one before launching...\n\n"
		confPicker := m.configPicker.View()
		content = header + confPicker
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
