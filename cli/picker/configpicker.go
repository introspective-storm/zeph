package picker

import tea "charm.land/bubbletea/v2"

//note that this is a case out of the pickers where we actually use it as a final result
// so, while the others are pure utlities to be used elsewhere before wrapped up in a
// nice little package and being sent to main, this one *also* needs a nice wrap up,
// on top of the basic raw utlity.

type ConfigLoader struct {
	picker Model
	Done   bool
	Path   string
}

func ConfigPicker(height int) Model {
	return New([]string{".zeph.yaml", ".zeph.yml"}, height)
}

func NewConfigLoader(height int) ConfigLoader {
	return ConfigLoader{
		picker: ConfigPicker(height),
	}
}

func (c ConfigLoader) Update(msg tea.Msg) (ConfigLoader, tea.Cmd) {
	var cmd tea.Cmd
	c.picker, cmd = c.picker.Update(msg)

	if selected, path := c.picker.CheckSelection(msg); selected {
		c.Path = path
		c.Done = true
	}
	return c, cmd
}

func (c ConfigLoader) View() string {
	return c.picker.View()
}

func (c ConfigLoader) Init() tea.Cmd {
	return c.picker.Init()
}
