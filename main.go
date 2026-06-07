package main

import (
	"netwatch/inputs"

	tea "charm.land/bubbletea/v2"
)

func main() {
	_, err := tea.NewProgram(inputs.InitModel()).Run()
	if err != nil {
		panic(err)
	}
}
