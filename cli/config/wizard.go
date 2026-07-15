package config

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"cli/picker"
)

type WizardStep int

const (
	StepProjectName  WizardStep = iota // Step 1: Project Name
	StepDatabasePath                   // Step 2: DB Picker
	StepModelPath                      // Step 3: Model Picker
	StepTests                          // Step 4: Choose Tests to Run
	StepConfirm                        // Step 5: Confirm and Save
)

// TestOption represents a diagnostic check the user can run.
type TestOption struct {
	Name        string
	Description string
	Selected    bool
}

type Wizard struct {
	step        WizardStep
	projectName string
	dbPath      string
	modelPath   string

	// Multi-select test suite
	tests         []TestOption
	activeTestIdx int

	nameInput   textinput.Model
	dbPicker    picker.Model
	modelPicker picker.Model

	Done bool
	Err  error
}

func NewWizard(height int) Wizard {
	ti := textinput.New()
	ti.Placeholder = "my-project"
	ti.Focus()
	ti.CharLimit = 64
	ti.SetWidth(40)

	// Standard Model-Agnostic Diagnostics & Test Suite
	availableTests := []TestOption{
		{
			Name:        "Default",
			Description: "Standard Suite of common tests",
			Selected:    true,
		},
		{
			Name:        "R-Squared (R²)",
			Description: "Measure variance explained by the model's regression targets",
			Selected:    true,
		},
		{
			Name:        "Confusion Matrix",
			Description: "Compute discrete classification matrices (TP, FP, TN, FN) for interval forecasting",
			Selected:    false,
		},
		{
			Name:        "MSE & RMSE",
			Description: "Track Mean Squared Error and Root Mean Squared Error to quantify error variance",
			Selected:    true,
		},
		{
			Name:        "Precision & Recall",
			Description: "Evaluate exact positive predictive value and sensitivity across threshold targets",
			Selected:    false,
		},
		{
			Name:        "Residual Plot Data",
			Description: "Output diagnostic residuals (y - y_hat) to check for underlying heteroscedasticity",
			Selected:    false,
		},
	}

	return Wizard{
		step:        StepProjectName,
		nameInput:   ti,
		dbPicker:    picker.DataPicker(height),
		modelPicker: picker.ModelPicker(height),
		tests:       availableTests,
	}
}

func (w Wizard) Update(msg tea.Msg) (Wizard, tea.Cmd) {
	var cmd tea.Cmd

	switch w.step {
	case StepProjectName:
		w.nameInput, cmd = w.nameInput.Update(msg)

		if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "enter" {
			name := strings.TrimSpace(w.nameInput.Value())
			if name == "" {
				name = w.nameInput.Placeholder
			}
			w.projectName = name
			w.step = StepDatabasePath
		}
		return w, cmd

	case StepDatabasePath:
		w.dbPicker, cmd = w.dbPicker.Update(msg)

		if selected, path := w.dbPicker.CheckSelection(msg); selected {
			w.dbPath = path
			w.step = StepModelPath
		}
		return w, cmd

	case StepModelPath:
		w.modelPicker, cmd = w.modelPicker.Update(msg)

		if selected, path := w.modelPicker.CheckSelection(msg); selected {
			w.modelPath = path
			w.step = StepTests // Step into multi-select suite
		}
		return w, cmd

	case StepTests:
		if key, ok := msg.(tea.KeyPressMsg); ok {
			switch key.String() {
			case "up", "k":
				if w.activeTestIdx > 0 {
					w.activeTestIdx--
				}
			case "down", "j":
				if w.activeTestIdx < len(w.tests)-1 {
					w.activeTestIdx++
				}
			case " ", "enter":
				// Toggle selection state
				w.tests[w.activeTestIdx].Selected = !w.tests[w.activeTestIdx].Selected
			case "y", "tab":
				// Proceed to final review
				w.step = StepConfirm
			}
		}
		return w, nil

	case StepConfirm:
		if key, ok := msg.(tea.KeyPressMsg); ok {
			switch key.String() {
			case "y", "enter":
				// Accumulate selected diagnostic configurations
				var selectedNames []string
				for _, t := range w.tests {
					if t.Selected {
						selectedNames = append(selectedNames, t.Name)
					}
				}

				// Clean config structure utilizing the renamed Tests slice!
				conf := &ProjectConfig{
					Name:       w.projectName,
					CreatedAt:  time.Now(),
					LastOpened: time.Now(),
					DataSource: w.dbPath,
					ModelPath:  w.modelPath,
					Tests:      selectedNames,
				}

				fileName := fmt.Sprintf("%s.tahoe.yaml", w.projectName)
				if err := conf.Save(fileName); err != nil {
					w.Err = err
					return w, nil
				}

				w.Done = true
			case "n":
				return NewWizard(15), nil
			}
		}
	}

	return w, nil
}

func (w Wizard) View() string {
	var s strings.Builder

	switch w.step {
	case StepProjectName:
		s.WriteString("Step 1: Enter Project Name\n\n")
		s.WriteString(w.nameInput.View())

	case StepDatabasePath:
		s.WriteString(fmt.Sprintf("Project Name: %s\n\n", w.projectName))
		s.WriteString("Step 2: Select Data Source File\n")
		s.WriteString(w.dbPicker.View())

	case StepModelPath:
		s.WriteString(fmt.Sprintf("Project Name: %s\n", w.projectName))
		s.WriteString(fmt.Sprintf("Data Source:  %s\n\n", w.dbPath))
		s.WriteString("Step 3: Select Model File\n")
		s.WriteString(w.modelPicker.View())

	case StepTests:
		s.WriteString(fmt.Sprintf("Project Name: %s\n", w.projectName))
		s.WriteString("Step 4: Select Diagnostic Diagnostics & Tests\n")
		s.WriteString("Choose what evaluations Tahoe will run against model runs:\n\n")

		for i, test := range w.tests {
			cursor := " "
			if i == w.activeTestIdx {
				cursor = "❯" // Highlighting active cursor line
			}

			checked := " "
			if test.Selected {
				checked = "x" // Selected visual indicator
			}

			// Style the active item line differently to draw attention to it
			if i == w.activeTestIdx {
				s.WriteString(fmt.Sprintf("  %s [%s] \033[1;36m%-22s\033[0m\n", cursor, checked, test.Name))
				s.WriteString(fmt.Sprintf("        \033[3m\033[90m%s\033[0m\n\n", test.Description))
			} else {
				s.WriteString(fmt.Sprintf("  %s [%s] %-22s\n", cursor, checked, test.Name))
				s.WriteString(fmt.Sprintf("        \033[90m%s\033[0m\n\n", test.Description))
			}
		}

		s.WriteString("(use ↑/↓ or j/k to navigate, Space/Enter to toggle, 'y' or Tab to continue)")

	case StepConfirm:
		s.WriteString("Step 5: Confirm Setup Details\n\n")
		s.WriteString(fmt.Sprintf("  Project Name: %s\n", w.projectName))
		s.WriteString(fmt.Sprintf("  Data Source:  %s\n", w.dbPath))
		s.WriteString(fmt.Sprintf("  Model File:   %s\n", w.modelPath))

		var activeTests []string
		for _, t := range w.tests {
			if t.Selected {
				activeTests = append(activeTests, t.Name)
			}
		}
		s.WriteString(fmt.Sprintf("  Enabled Tests: %s\n\n", strings.Join(activeTests, ", ")))
		s.WriteString(fmt.Sprintf("Create config file (%s.tahoe.yaml)? (y/n)", w.projectName))
	}

	return s.String()
}
