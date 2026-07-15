package config

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func selectModel(msg tea.Msg, w Wizard) (Wizard, tea.Cmd) {
	var cmd tea.Cmd
	w.modelPicker, cmd = w.modelPicker.Update(msg)

	if selected, path := w.modelPicker.CheckSelection(msg); selected {
		w.modelPath = path
		w.step = StepTests // Step into multi-select suite
	}
	return w, cmd
}

func renderModel(w Wizard) string {
	var s strings.Builder
	s.WriteString("3/4\n")
	s.WriteString(fmt.Sprintf("Project Name: %s\n", w.projectName))
	s.WriteString(fmt.Sprintf("Data Source:  %s\n\n", w.dbPath))
	s.WriteString(" Select Model File\n")
	s.WriteString(w.modelPicker.View())
	return s.String()
}
