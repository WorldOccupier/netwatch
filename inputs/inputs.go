package inputs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	inptusCount    = 2
	inputCharLimit = 100
	numbersOnlyRegexMatcher = "^[0-9]+$"
	quitShortcut   = "ctrl+c"
	esc 		   = "esc"
	tab            = "tab"
	reverseTab     = "shift+tab"
	enter          = "enter"
	up             = "up"
	down           = "down"
	backspace      = "backspace"

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type model struct {
	focusIndex int
	inputs []textinput.Model
	quitting bool
}

func InitModel() model {
	m := model {
		inputs: make([]textinput.Model, inptusCount),
	}

	for i := range m.inputs {
		m.inputs[i] = m.getDefaultTextInput(i)
	}

	return m
}

func (m model) getDefaultTextInput(i int) textinput.Model {
	textInput :=	textinput.New()
	textInput.CharLimit = inputCharLimit
	textInput.SetWidth(inputCharLimit)

	style := textInput.Styles()
	style.Cursor.Color = lipgloss.Color("205")
	style.Focused.Prompt = focusedStyle
	style.Focused.Text = focusedStyle
	style.Blurred.Prompt = blurredStyle
	style.Focused.Text = focusedStyle
	textInput.SetStyles(style)

	switch i {
	case 0:
		textInput.Placeholder = "Stocks and Shares"
		textInput.Focus()
	case 1:
		textInput.Placeholder = "Savings"
	}

	return textInput
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case quitShortcut, esc:
			m.quitting = true
			return m, tea.Quit
		case tab, reverseTab, enter, up, down:
			if msg.String() == enter && m.focusIndex == len(m.inputs) {
				m.quitting = true
				return m, tea.Batch()
			}
			cmds := m.handleMovement(&msg)
			return m, tea.Batch(cmds...)
		}

		m.updateInputs(msg)
	}

	return m, tea.Batch()
}

func (m *model) handleMovement(msg *tea.KeyPressMsg) []tea.Cmd {
	m.setFocusIndex(msg)
	cmds := m.handleInputsFocus()
	return cmds
}

func (m *model) setFocusIndex(msg *tea.KeyPressMsg) {
	msgString := msg.String()
	switch msgString {
	case up, reverseTab:
		m.focusIndex--
	case down, tab, enter:
		m.focusIndex++
	}

	if m.focusIndex > len(m.inputs) {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs)
	}
}

func (m *model) handleInputsFocus() []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
			continue
		}

		m.inputs[i].Blur()
	}

	return cmds
}

func (m model) updateInputs(msg tea.KeyPressMsg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		shouldUpdate := shouldUpdateInputs(msg)
		if shouldUpdate {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		}
	}

	return tea.Batch(cmds...)
}

func shouldUpdateInputs(msg tea.KeyPressMsg) bool {
	var inputValidationRegex = regexp.MustCompile(numbersOnlyRegexMatcher)
	inputValidated := inputValidationRegex.MatchString(msg.String())
	msgString := msg.String()
	return inputValidated || msgString == backspace || msgString == "left" || msgString == "right"
}

func (m model) View() tea.View {
	if m.quitting {
		return m.resultsView()
	}

	return m.inputsView()
}

func (m model) resultsView() tea.View {
	var total int64
	for _, textInput := range m.inputs {
		stringValue := textInput.Value()
		if stringValue != "" {
			intValue, err := strconv.ParseInt(stringValue, 10, 64)
			if err != nil {
				panic(err)
			}
			total += intValue
		}
	}
	var stringBuilder strings.Builder
	fmt.Fprintf(&stringBuilder, "%s", blurredStyle.Render("Total"))
	stringBuilder.WriteRune('\n')
	fmt.Fprintf(&stringBuilder, "%s", focusedStyle.Render(strconv.FormatInt(total, 10)))
	stringBuilder.WriteRune('\n')
	return tea.NewView(stringBuilder.String())
}

func (m model) inputsView() tea.View {
	var stringBuilder strings.Builder
	for i := range m.inputs {
		stringBuilder.WriteString(m.inputs[i].View())
		stringBuilder.WriteRune('\n')
	}

	focusedButton := lipgloss.NewStyle().Foreground(lipgloss.Color("#0cce2c")).Render("[ Submit ]")
	blurredButton := fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}

	fmt.Fprintf(&stringBuilder, "\n%s\n", *button)

	return tea.NewView(stringBuilder.String())
}
