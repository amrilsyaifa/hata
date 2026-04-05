package locale

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// Select shows an interactive multi-select list of all locales.
// The user navigates with arrow keys, toggles with Space, and confirms with Enter.
// Typing narrows the list instantly.
func Select() ([]string, error) {
	options := make([]string, len(All))
	for i, l := range All {
		options[i] = l.Display()
	}

	var selected []string
	prompt := &survey.MultiSelect{
		Message:  "Select languages (Space to toggle, Enter to confirm, type to filter):",
		Options:  options,
		PageSize: 15,
	}
	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, fmt.Errorf("language selection cancelled: %w", err)
	}

	codes := make([]string, 0, len(selected))
	for _, s := range selected {
		codes = append(codes, CodeFromDisplay(s))
	}
	return codes, nil
}
