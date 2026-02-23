package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#33bf24"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240"))

	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212"))

	buttonStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("99")).
			Foreground(lipgloss.Color("0")).
			Padding(0, 2).
			Bold(true)
)

func panelStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1).
		Width(width).
		Height(height)
}

type (
	Method int
	Mode   int
)

const (
	Caesar Method = iota
	Vigenere
)

const (
	Encrypt Mode = iota
	Decrypt
)

const (
	FocusMethod int = iota
	FocusMode
	FocusKey
	FocusShift
	FocusInput
	FocusOutput
	FocusBtn
	FocusCount
)

type model struct {
	Method Method
	Mode   Mode
	Focus  int
	Width  int
	Height int
	Key    textinput.Model
	Shift  textinput.Model
	Input  textarea.Model
	Output textarea.Model
	keymap keymap
}

type keymap = struct {
	next, prev, quit key.Binding
}

func newTextarea(placeholder string) textarea.Model {
	ta := textarea.New()
	ta.Cursor.Style = cursorStyle
	ta.Prompt = ""
	ta.Placeholder = placeholder
	ta.ShowLineNumbers = false
	ta.FocusedStyle.Base = focusedBorderStyle
	ta.BlurredStyle.Base = blurredBorderStyle
	ta.Focus()

	return ta
}

func newOutputTextarea(placeholder string) textarea.Model {
	ta := textarea.New()
	ta.Prompt = ""
	ta.Placeholder = placeholder
	ta.ShowLineNumbers = false
	ta.FocusedStyle.Base = focusedBorderStyle
	ta.BlurredStyle.Base = blurredBorderStyle

	return ta
}

func newTextinput(placeholder string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Cursor.Style = cursorStyle
	return ti
}

func NewModel() model {
	m := model{
		Method: Caesar,
		Mode:   Encrypt,
		Focus:  FocusInput,
		Input:  newTextarea("Type text you want to cipher..."),
		Output: newOutputTextarea("Result will appear here..."),
		Key:    newTextinput("Enter key..."),
		Shift:  newTextinput("3"),
		keymap: keymap{
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
			),
			next: key.NewBinding(
				key.WithKeys("tab"),
			),
			prev: key.NewBinding(
				key.WithKeys("shift+tab"),
			),
		},
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		leftWidth := 25
		totalH := msg.Height - 2
		topContentH := (totalH - 8) / 2
		botContentH := (totalH-8)/2 + (totalH-8)%2
		rightContentW := msg.Width - leftWidth - 8 // -4 left panel (border+pad), -4 right panel (border+pad)

		// textarea overhead inside panel: 1 (label) + 2 (textarea border) = 3
		m.Input.SetWidth(rightContentW - 2) // -2 for textarea border left+right
		m.Input.SetHeight(topContentH - 3)
		m.Output.SetWidth(rightContentW - 2)
		m.Output.SetHeight(botContentH - 3)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.Input.Blur()
			m.Output.Blur()
			m.Key.Blur()
			m.Shift.Blur()
			return m, tea.Quit
		case key.Matches(msg, m.keymap.next):
			m.Focus = (m.Focus + 1) % FocusCount
			m.updateFocus()
			return m, nil
		case key.Matches(msg, m.keymap.prev):
			m.Focus = (m.Focus - 1 + FocusCount) % FocusCount
			m.updateFocus()
			return m, nil
		case m.Focus == FocusMethod && msg.String() == "enter":
			m.Method = (m.Method + 1) % 2
			return m, nil
		case m.Focus == FocusMode && msg.String() == "enter":
			m.Mode = (m.Mode + 1) % 2
			return m, nil
		}
	}

	var cmd tea.Cmd

	if m.Focus == FocusInput {
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.Focus == FocusOutput {
		m.Output, cmd = m.Output.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.Focus == FocusKey {
		m.Key, cmd = m.Key.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.Focus == FocusShift {
		m.Shift, cmd = m.Shift.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) updateFocus() {
	m.Input.Blur()
	m.Key.Blur()
	m.Shift.Blur()

	switch m.Focus {
	case FocusInput:
		m.Input.Focus()
	case FocusKey:
		m.Key.Focus()
	case FocusShift:
		m.Shift.Focus()
	}
}

func (m model) View() string {
	methodLabel := labelStyle.Render("METHOD")
	methodValue := "Caesar"
	if m.Method == Vigenere {
		methodValue = "Vigenere"
	}
	if m.Focus == FocusMethod {
		methodValue = selectedStyle.Render(methodValue) + " (Enter)"
	}

	modeLabel := labelStyle.Render("MODE")
	modeValue := "Encrypt"
	if m.Mode == Decrypt {
		modeValue = "Decrypt"
	}
	if m.Focus == FocusMode {
		modeValue = selectedStyle.Render(modeValue) + " (Enter)"
	}

	keyLabel := labelStyle.Render("KEY")
	shiftLabel := labelStyle.Render("SHIFT")

	runBtnLabel := "RUN"
	if m.Focus == FocusBtn {
		runBtnLabel = buttonStyle.Render(runBtnLabel)
	}

	leftPanel := lipgloss.JoinVertical(
		lipgloss.Left,
		methodLabel,
		methodValue,
		"",
		modeLabel,
		modeValue,
		"",
		keyLabel,
		m.Key.View(),
		"",
		shiftLabel,
		m.Shift.View(),
		"",
		"",
		runBtnLabel,
	)

	leftWidth := 25
	rightWidth := m.Width - leftWidth - 8 // leftWidth + 4 (left borders+padding) + 4 (right borders+padding)

	totalH := m.Height - 2
	rem := (totalH - 8) % 2
	topH := (totalH - 8) / 2
	botH := topH + rem
	leftH := totalH - 4 // всегда равно topOuterH + botOuterH = totalH

	leftPanel = panelStyle(leftWidth, leftH).Render(leftPanel)

	inputPanelStyle := lipgloss.NewStyle().
		Height(topH).
		Width(rightWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1)

	outputPanelStyle := lipgloss.NewStyle().
		Height(botH).
		Width(rightWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1)

	inputPanel := inputPanelStyle.Render(labelStyle.Render("INPUT") + "\n" + m.Input.View())

	outputPanel := outputPanelStyle.Render(labelStyle.Render("OUTPUT") + "\n" + m.Output.View())

	rightPanel := lipgloss.JoinVertical(0, inputPanel, outputPanel)

	return lipgloss.JoinHorizontal(0, leftPanel, rightPanel)
}
