package main

import tea "charm.land/bubbletea/v2"

type AppState int

const (
	StateMenu AppState = iota
	StateWizard
	StateConfigPicker
	StateDataPicker
	StateModelPicker
)

func (m rootModel) stateMenu(msg tea.Msg) (rootModel, tea.Cmd) {
	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)

	if keyMsg, ok := msg.(tea.KeyPressMsg); ok && keyMsg.String() == "enter" {
		if selected, ok := m.menu.SelectedItem().(item); ok {
			switch selected.GetAction() {
			case "wizard":
				m.state = StateWizard
			case "load":
				m.state = StateConfigPicker
			}
		}
	}
	return m, cmd
}
