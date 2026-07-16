package config

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

func selectSaveLocation(msg tea.Msg, w Wizard) (Wizard, tea.Cmd) {
	var cmd tea.Cmd
	w.dirPicker, cmd = w.dirPicker.Update(msg)

	if selected, path := w.dirPicker.CheckSelection(msg); selected {
		w.saveDir = path
		w.step = StepConfirm
	}
	return w, cmd
}

func renderSaveLocation(w Wizard) string {
	var s strings.Builder
	s.WriteString("5/5\n")
	s.WriteString("Save Location\n")
	s.WriteString("Select the directory where you'd like to save your configuration:\n\n")
	s.WriteString(w.dirPicker.View())
	return s.String()
}
