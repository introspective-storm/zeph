package config

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type ProjectStep struct {
	input textinput.Model
}

func newProjectInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "my-project"
	ti.Focus()
	ti.CharLimit = 64
	ti.SetWidth(40)
	return ti
}

func nameProject(msg tea.Msg, w Wizard) (Wizard, tea.Cmd) {
	var cmd tea.Cmd
	w.nameInput, cmd = w.nameInput.Update(msg)

	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "enter" {
		name := strings.TrimSpace(w.nameInput.Value())
		if name == "" {
			name = w.nameInput.Placeholder
		}
		w.projectName = name
		w.step = StepDatabasePath
		return w, w.dbPicker.Init()
	}
	return w, cmd
}

func renderProject(w Wizard) string {
	var s strings.Builder
	s.WriteString("1/5\n")
	s.WriteString("Enter Project Name\n\n")
	s.WriteString(w.nameInput.View())
	return s.String()
}
