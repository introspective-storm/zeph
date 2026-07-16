package config

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"cli/picker"
)

type WizardStep int

const (
	StepProjectName WizardStep = iota
	StepDatabasePath
	StepModelPath
	StepTests
	StepConfirm
)

type Wizard struct {
	step        WizardStep
	projectName string
	dbPath      string
	modelPath   string

	tests         []TestOption
	activeTestIdx int

	nameInput   textinput.Model
	dbPicker    picker.Model
	modelPicker picker.Model

	Done bool
	Err  error
}

func NewWizard(height int) Wizard {

	return Wizard{
		step:        StepProjectName,
		nameInput:   newProjectInput(),
		dbPicker:    picker.DataPicker(height),
		modelPicker: picker.ModelPicker(height),
		tests:       GetAvailableTests(),
	}
}

func (w Wizard) Update(msg tea.Msg) (Wizard, tea.Cmd) {
	switch w.step {
	case StepProjectName:
		return nameProject(msg, w)
	case StepDatabasePath:
		return selectData(msg, w)
	case StepModelPath:
		return selectModel(msg, w)
	case StepTests:
		return selectTest(msg, w)
	case StepConfirm:
		return confirmProject(msg, w)
	}
	return w, nil
}

func (w Wizard) View() string {
	var s strings.Builder

	switch w.step {
	case StepProjectName:
		s.WriteString(renderProject(w))
	case StepDatabasePath:
		s.WriteString(renderData(w))
	case StepModelPath:
		s.WriteString(renderModel(w))
	case StepTests:
		s.WriteString(renderTest(w))
	case StepConfirm:
		s.WriteString(renderComfirm(w))
	}
	return s.String()
}

func (w Wizard) Init() tea.Cmd {
	return w.nameInput.Focus()
}
