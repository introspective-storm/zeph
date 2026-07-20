package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

func confirmProject(msg tea.Msg, w Wizard) (Wizard, tea.Cmd) {
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

			newConf := &ProjectConfig{
				Name:       w.projectName,
				CreatedAt:  time.Now(),
				LastOpened: time.Now(),
				DataSource: w.dbPath,
				ModelPath:  w.modelPath,
				Tests:      selectedNames,
			}

			fileName := fmt.Sprintf("%s.zeph.yaml", w.projectName)
			absoluteFilePath := filepath.Join(w.saveDir, fileName)

			if err := newConf.Save(absoluteFilePath); err != nil {
				w.Err = err
				return w, nil
			}
			w.Conf = newConf
			w.Done = true
		case "n":
			return NewWizard(15), nil
		}
	}
	return w, nil
}

func renderComfirm(w Wizard) string {
	var s strings.Builder
	s.WriteString("Confirm Setup Details\n\n")
	s.WriteString(fmt.Sprintf("  Project Name: %s\n", w.projectName))
	s.WriteString(fmt.Sprintf("  Data Source:  %s\n", w.dbPath))
	s.WriteString(fmt.Sprintf("  Model File:   %s\n", w.modelPath))
	s.WriteString(fmt.Sprintf("  Save Folder:   %s\n", w.saveDir))

	var activeTests []string
	for _, t := range w.tests {
		if t.Selected {
			activeTests = append(activeTests, t.Name)
		}
	}
	s.WriteString(fmt.Sprintf("  Enabled Tests: %s\n\n", strings.Join(activeTests, ", ")))
	s.WriteString(fmt.Sprintf("Create config file (%s.zeph.yaml)? [Y/n]", w.projectName))
	return s.String()
}
