package picker

import (
	"os"

	"charm.land/bubbles/v2/filepicker"
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	picker filepicker.Model
}

func New(allowedExt []string, height int) Model {
	p := filepicker.New()
	p.AllowedTypes = allowedExt
	p.SetHeight(height)

	if dir, err := os.Getwd(); err == nil {
		p.CurrentDirectory = dir
	}

	return Model{
		picker: p,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.picker, cmd = m.picker.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.picker.View()
}

func (m Model) Init() tea.Cmd {
	return m.picker.Init()
}

func (m *Model) CheckSelection(msg tea.Msg) (bool, string) {
	return m.picker.DidSelectFile(msg)
}
