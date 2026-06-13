package models

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type currency struct {
	values []string
	focusIndex int
	isFocused bool
}


func (c currency) Init() tea.Cmd {
	return nil
}

func (c currency) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyPressMsg:
		switch msg.String() {
		case quitShortcut, esc:
			return c, tea.Quit
		case left:
			c.focusIndex--
			if c.focusIndex < 0 {
				c.focusIndex = len(c.values) - 1
			}
		case right:
			c.focusIndex++
			if c.focusIndex >= len(c.values) {
				c.focusIndex = 0
			}
		}
	}

	return c, tea.Batch()
}

func (c currency) View() tea.View {
	var stringBuilder strings.Builder
	var style lipgloss.Style
	if c.isFocused {
		style = focusedStyle
	} else {
		style = blurredStyle
	}
	fmt.Fprintf(&stringBuilder, "%s", style.Render(c.values[c.focusIndex]))

	return tea.NewView(stringBuilder.String())
}

func (c currency) GetActiveValue() string {
	return c.values[c.focusIndex]
}

func (c currency) SetActiveIndex(value string) {
	for i := range c.values {
		if c.values[i] == value {
			c.focusIndex = i
			break
		}
	}
}
