package models

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type display struct {
	inputsArray []model
}

func InitDisplay() display {
	return display {
		inputsArray: []model { InitModel() },
	}
}

func (d display) Init() tea.Cmd {
	return textinput.Blink
}

func (d display) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return d.inputsArray[0].Update(msg)
}

func (d display) View() tea.View {
	return d.inputsArray[0].View()
}
