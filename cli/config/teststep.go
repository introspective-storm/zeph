package config

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type TestOption struct {
	Name        string
	Description string
	Selected    bool
}

func GetAvailableTests() []TestOption {
	return []TestOption{
		{
			Name:        "Default",
			Description: "Standard Suite of common tests",
			Selected:    true,
		},
		{
			Name:        "MAE",
			Description: "Track the Mean Average Error to quantify the error of the model's predicted targets",
			Selected:    false,
		},
		{
			Name:        "R-Squared",
			Description: "Measure variance explained by the model's regression targets",
			Selected:    false,
		},
		{
			Name:        "MSE & RMSE",
			Description: "Track Mean Squared Error and Root Mean Squared Error to quantify error variance",
			Selected:    false,
		},
		{
			Name:        "Confusion Matrix",
			Description: "Compute discrete classification matrices (TP, FP, TN, FN) for interval forecasting",
			Selected:    false,
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
		{
			Name:        "Permutation Importance",
			Description: "Assess feature impact by shuffling input columns to measure performance degradation",
			Selected:    false,
		},
		{
			Name:        "Calibration Curve",
			Description: "Map predicted vs. actual outcome frequencies to identify bias/overconfidence",
			Selected:    false,
		},
		{
			Name:        "Population Stability Index (PSI)",
			Description: "Quantify the shift in feature distributions between training and live inference",
			Selected:    false,
		},
		{
			Name:        "Ljung-Box Test",
			Description: "Check for residual autocorrelation to ensure all signal has been extracted",
			Selected:    false,
		},
		{
			Name:        "Brier Score",
			Description: "Measure the mean squared difference between predicted probability and actual outcome",
			Selected:    false,
		},
	}
}

func selectTest(msg tea.Msg, w Wizard) (Wizard, tea.Cmd) {
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
		case "x", "tab":
			w.tests[w.activeTestIdx].Selected = !w.tests[w.activeTestIdx].Selected
		case "y", "enter":
			w.step = StepConfirm
		}
	}
	return w, nil
}

func renderTest(w Wizard) string {
	var s strings.Builder
	s.WriteString("4/4\n")
	s.WriteString(fmt.Sprintf("Project Name: %s\n", w.projectName))
	s.WriteString("Select Diagnostic Diagnostics & Tests\n")
	s.WriteString("Choose what evaluations Zeph will run against your models:\n\n")

	for i, test := range w.tests {
		cursor := " "
		if i == w.activeTestIdx {
			cursor = "❯" //active cursor line
		}

		checked := " "
		if test.Selected {
			checked = "x" //visual indicator for 'selected'
		}

		// styling for active and unactive are different
		if i == w.activeTestIdx {
			s.WriteString(fmt.Sprintf("  %s [%s] \033[1;36m%-22s\033[0m\n", cursor, checked, test.Name))
			s.WriteString(fmt.Sprintf("        \033[3m\033[90m%s\033[0m\n\n", test.Description))
		} else {
			s.WriteString(fmt.Sprintf("  %s [%s] %-22s\n", cursor, checked, test.Name))
			s.WriteString(fmt.Sprintf("        \033[90m%s\033[0m\n\n", test.Description))
		}
	}

	s.WriteString("(use ↑/↓ or j/k to navigate, 'x' or Tab to toggle, 'y' or Enter to continue)")

	return s.String()
}
