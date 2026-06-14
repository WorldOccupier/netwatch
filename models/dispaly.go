package models

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type display struct {
	inputsArray []model
	focusIndex int
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
	var cmd tea.Cmd
	var teaModel tea.Model
	for i := range d.inputsArray {
		teaModel, cmd = d.inputsArray[i].Update(msg)
		d.inputsArray[i] = teaModel.(model)
	}
	return d, cmd
}

func (d display) View() tea.View {
	inputsView := d.inputsArray[0].View()
	// focusedButton := lipgloss.NewStyle().Foreground(lipgloss.Color("#0cce2c")).Render("[ Submit ]")
	blurredButton := fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	var stringBuilder strings.Builder
	stringBuilder.WriteString(inputsView.Content)

	button := &blurredButton
	// if m.focusIndex == m.count - 1 {
	// 	button = &focusedButton
	// }

	fmt.Fprintf(&stringBuilder, "\n%s\n", *button)

	return tea.NewView(stringBuilder.String())
}
