package main

import (
	"netwatch/models"

	tea "charm.land/bubbletea/v2"
)

func main() {
	_, err := tea.NewProgram(models.InitModel()).Run()
	if err != nil {
		panic(err)
	}
}
