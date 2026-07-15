package picker

func ConfigPicker(height int) Model {
	return New([]string{".tahoe.yaml", ".tahoe.yml"}, height)
}
