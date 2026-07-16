package picker

import (
	"os"

	"charm.land/bubbles/v2/filepicker"
)

func DirectoryPicker(height int) Model {
	p := filepicker.New()
	p.DirAllowed = true
	p.FileAllowed = false
	p.SetHeight(height)

	if dir, err := os.Getwd(); err == nil {
		p.CurrentDirectory = dir
	}
	return Model{
		picker: p,
	}
}
