package config

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func selectData(msg tea.Msg, w Wizard) (Wizard, tea.Cmd) {
	var cmd tea.Cmd
	w.dbPicker, cmd = w.dbPicker.Update(msg)

	if selected, path := w.dbPicker.CheckSelection(msg); selected {
		w.dbPath = path
		w.step = StepModelPath
		return w, w.modelPicker.Init()
	}
	return w, cmd
}

func renderData(w Wizard) string {
	var s strings.Builder
	s.WriteString("2/5\n")
	s.WriteString(fmt.Sprintf("Project Name: %s\n\n", w.projectName))
	s.WriteString("Select Data Source File\n")
	s.WriteString(w.dbPicker.View())
	return s.String()
}
