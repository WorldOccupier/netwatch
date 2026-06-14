package models

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
	left 		   = "left"
	right    	   = "right"
	backspace      = "backspace"
	currencyDisplayPosition = 0
	focusedColor = "205"

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(focusedColor))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type model struct {
	focusIndex int
	currency currency
	inputs []textinput.Model
	quitting bool
	submitted bool
	count int
	currencyIndex int
}

func InitModel() model {
	savedCurrency, savedInputs := loadLatestSavedDate()
	currencyValues := []string {"GBP", "INR"}
	var savedCurrencyIndex int
	for i := range currencyValues {
		if currencyValues[i] == savedCurrency {
			savedCurrencyIndex = i
			break
		}
	}
	inputCurrency := currency {values: currencyValues, focusIndex: savedCurrencyIndex}
	m := model {
		currency: inputCurrency,
		inputs: make([]textinput.Model, inptusCount),
		count: 4,
	}
	m.currency.isFocused = true

	for i := range m.inputs {
		m.inputs[i] = m.getDefaultTextInput(i)
		m.inputs[i].SetValue(savedInputs[i])
	}

	return m
}

func (m model) getDefaultTextInput(i int) textinput.Model {
	textInput :=	textinput.New()
	textInput.CharLimit = inputCharLimit
	textInput.SetWidth(inputCharLimit)

	style := getInputStyle(textInput)
	textInput.SetStyles(style)
	inputPlaceholder := "$$$"
	textInput.Placeholder = inputPlaceholder

	switch i {
	case 0:
		textInput.Prompt = "Stocks and Shares: "
	case 1:
		textInput.Prompt = "Savings:           "
	}

	return textInput
}

func getInputStyle(textInput textinput.Model) textinput.Styles {
	style := textInput.Styles()
	style.Cursor.Color = lipgloss.Color(focusedColor)
	style.Focused.Prompt = focusedStyle
	style.Focused.Text = focusedStyle
	style.Blurred.Prompt = blurredStyle
	style.Focused.Text = focusedStyle

	return style
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
			if msg.String() == enter {
				if m.submitted {
					return m, tea.Quit
				}
				if m.focusIndex == m.count - 1 {
					m.submitted = true
					m.save()
					return m, tea.Batch()
				}
			}
			cmds := m.handleMovement(&msg)
			return m, tea.Batch(cmds...)
		case left, right:
			if m.focusIndex == currencyDisplayPosition {
				currencyModel, cmd := m.currency.Update(msg)
				m.currency = currencyModel.(currency)
				return m, cmd
			}
		}

		return m, m.updateInputs(msg)
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

	if m.focusIndex >= m.count {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = m.count - 1
	}
	if m.focusIndex == currencyDisplayPosition {
		m.currency.isFocused = true
	} else {
		m.currency.isFocused = false
	}
}

func (m *model) handleInputsFocus() []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range len(m.inputs) {
		if i == m.focusIndex - 1  {
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
	if m.submitted {
		return m.resultsView()
	}

	return m.inputsView()
}

func (m model) resultsView() tea.View {
	total := m.total()
	var stringBuilder strings.Builder
	fmt.Fprintf(&stringBuilder, "%s", blurredStyle.Render("Total"))
	stringBuilder.WriteRune('\n')
	fmt.Fprintf(&stringBuilder, "%s", focusedStyle.Render(strconv.FormatInt(total, 10)))
	stringBuilder.WriteRune('\n')
	return tea.NewView(stringBuilder.String())
}

func (m model) total() int64 {
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

	return total
}

func (m model) inputsView() tea.View {
	var stringBuilder strings.Builder

	currencyView := m.currency.View()
	stringBuilder.WriteString(currencyView.Content)
	stringBuilder.WriteRune('\n')
	for i := range m.inputs {
		stringBuilder.WriteString(m.inputs[i].View())
		stringBuilder.WriteRune('\n')
	}

	return tea.NewView(stringBuilder.String())
}
