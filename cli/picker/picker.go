package picker

import (
	"os"

	"charm.land/bubbles/v2/filepicker"
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type Model struct {
	picker filepicker.Model
	help   help.Model
}

type helpKeyMap struct {
	km filepicker.KeyMap
}

type helpKeyMapLocked struct {
	km filepicker.KeyMap
}

func (h helpKeyMap) ShortHelp() []key.Binding {
	up := h.km.Up
	down := h.km.Down
	open := h.km.Open
	selectBtn := h.km.Select
	back := h.km.Back

	up.SetHelp("↑/k", "up")
	down.SetHelp("↓/j", "down")
	open.SetHelp("enter", "open")
	selectBtn.SetHelp("enter", "select")
	back.SetHelp("backspace/h", "back")

	return []key.Binding{
		up,
		down,
		open,
		selectBtn,
		back,
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "cancel")),
		key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("qq/ctrl+c", "quit")), // Unified here
	}
}

func (h helpKeyMap) FullHelp() [][]key.Binding {
	up := h.km.Up
	down := h.km.Down
	pageUp := h.km.PageUp
	pageDown := h.km.PageDown
	open := h.km.Open
	selectBtn := h.km.Select
	back := h.km.Back
	goToTop := h.km.GoToTop
	goToLast := h.km.GoToLast

	up.SetHelp("↑/k", "up")
	down.SetHelp("↓/j", "down")
	open.SetHelp("enter", "open")
	selectBtn.SetHelp("enter", "select")
	back.SetHelp("backspace/h", "back")

	return [][]key.Binding{
		{up, down, pageUp, pageDown},
		{
			open,
			selectBtn,
			back,
			goToTop,
			goToLast,
			key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "cancel")),
			key.NewBinding(key.WithKeys("qq", "ctrl+c"), key.WithHelp("qq/ctrl+c", "quit")), // Unified here
		},
	}
}

func (h helpKeyMapLocked) ShortHelp() []key.Binding {
	up := h.km.Up
	down := h.km.Down
	selectBtn := h.km.Select

	up.SetHelp("↑/k", "up")
	down.SetHelp("↓/j", "down")
	selectBtn.SetHelp("enter", "select")

	return []key.Binding{
		up,
		down,
		selectBtn,
		key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
	}
}

func (h helpKeyMapLocked) FullHelp() [][]key.Binding {
	return [][]key.Binding{h.ShortHelp()}
}

// constructor for filepicker components
func New(allowedExt []string, height int) Model {
	p := filepicker.New()
	p.AllowedTypes = allowedExt
	p.SetHeight(height)

	if dir, err := os.Getwd(); err == nil {
		p.CurrentDirectory = dir
	}

	return Model{
		picker: p,
		help:   help.New(),
	}
}

func (m Model) pickerKeyMap() helpKeyMap {
	return helpKeyMap{km: m.picker.KeyMap}
}

func (m Model) pickerKeyMapLocked() helpKeyMapLocked {
	return helpKeyMapLocked{km: m.picker.KeyMap}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.picker, cmd = m.picker.Update(msg)
	return m, cmd
}

func (m Model) UpdateLocked(msg tea.Msg) (Model, tea.Cmd) {
	expectedDir, _ := os.Getwd()
	var cmd tea.Cmd
	m.picker, cmd = m.picker.Update(msg)
	if m.picker.CurrentDirectory != expectedDir {
		m.picker.CurrentDirectory = expectedDir
		return m, m.picker.Init()
	}
	return m, cmd
}

func (m Model) View() string {
	pickerView := m.picker.View()
	helpView := m.help.View(m.pickerKeyMap())
	return pickerView + helpView
}

func (m Model) ViewLocked() string {
	pickerView := m.picker.View()
	helpView := m.help.View(m.pickerKeyMapLocked())
	return pickerView + helpView
}

func (m Model) Init() tea.Cmd {
	return m.picker.Init()
}

func (m *Model) CheckSelection(msg tea.Msg) (bool, string) {
	return m.picker.DidSelectFile(msg)
}
