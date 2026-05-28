package ui

import (
	"os"

	"github.com/charmbracelet/huh"
)

func Confirm(message string, assumeYes bool) (bool, error) {
	if assumeYes {
		return true, nil
	}
	if os.Getenv("CI") != "" {
		return false, nil
	}
	var confirmed bool
	form := huh.NewForm(huh.NewGroup(huh.NewConfirm().Title(message).Value(&confirmed)))
	if err := form.Run(); err != nil {
		return false, err
	}
	return confirmed, nil
}
