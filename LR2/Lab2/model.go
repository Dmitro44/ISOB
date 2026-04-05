package main

import (
	"go-cipher/crypto"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Init() tea.Cmd {
	return m.Filepicker.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		layout := NewLayout(msg.Width, msg.Height)

		m.Input.SetWidth(layout.TextareaInputWidth)
		m.Input.SetHeight(layout.TextareaInputHeight)
		m.Output.SetWidth(layout.TextareaOutputWidth)
		m.Output.SetHeight(layout.TextareaOutputHeight)
		m.Filepicker.SetHeight(layout.FilepickerHeight)
		m.SaveFilename.Width = 34

	case statusMsg:
		m.Status = ""
		return m, nil

	case tea.KeyMsg:
		// Save Popup Logic
		if m.ShowSavePopup {
			switch msg.String() {
			case "esc":
				m.ShowSavePopup = false
				m.SaveFilename.Blur()
				return m, nil
			case "enter":
				filename := m.SaveFilename.Value()
				if filename == "" {
					filename = "output.txt"
				}

				err := os.WriteFile(filename, []byte(m.Output.Value()), 0o644)
				if err != nil {
					m.Status = "Error: " + err.Error()
				} else {
					m.Status = "Saved to: " + filename
				}

				m.ShowSavePopup = false
				m.SaveFilename.Blur()

				return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
					return statusMsg("")
				})
			default:
				var cmd tea.Cmd
				m.SaveFilename, cmd = m.SaveFilename.Update(msg)
				return m, cmd
			}
		}

		// Global Bindings
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.Input.Blur()
			m.Output.Blur()
			m.Key.Blur()
			m.Shift.Blur()
			return m, tea.Quit
		case key.Matches(msg, m.keymap.next):
			m.Focus = (m.Focus + 1) % FocusCount
			m.skipInactive(1)
			m.updateFocus()
			return m, m.Filepicker.Init()
		case key.Matches(msg, m.keymap.prev):
			m.Focus = (m.Focus - 1 + FocusCount) % FocusCount
			m.skipInactive(-1)
			m.updateFocus()
			return m, m.Filepicker.Init()
		}

		// Filepicker specific
		if m.InputMode == File && m.Focus == FocusFilepicker {
			if !m.FileConfirmed {
				var cmd tea.Cmd
				m.Filepicker, cmd = m.Filepicker.Update(msg)
				if didSelect, path := m.Filepicker.DidSelectFile(msg); didSelect {
					m.SelectedFile = path
					m.FileConfirmed = true
				}
				return m, cmd
			}
			if msg.String() == "backspace" {
				m.FileConfirmed = false
				m.SelectedFile = ""
				return m, m.Filepicker.Init()
			}
			return m, nil
		}

		// Other focusable elements
		switch {
		case key.Matches(msg, m.keymap.copy):
			m.Input.SetValue(m.Output.Value())
			m.Output.SetValue("")
			m.InputMode = ManualInput
			m.Focus = FocusInput
			m.updateFocus()
			return m, nil
		case key.Matches(msg, m.keymap.clearInput):
			m.Input.SetValue("")
			m.Focus = FocusInput
			m.updateFocus()
			return m, nil
		case m.Focus == FocusInputMode && msg.String() == "enter":
			m.InputMode = (m.InputMode + 1) % 2
			m.FileConfirmed = false
			m.SelectedFile = ""
			m.updateFocus()
			return m, m.Filepicker.Init()
		case m.Focus == FocusMethod && msg.String() == "enter":
			m.Method = (m.Method + 1) % 2
			m.skipInactive(1)
			m.updateFocus()
			return m, nil
		case m.Focus == FocusMode && msg.String() == "enter":
			m.Mode = (m.Mode + 1) % 2
			return m, nil
		case m.Focus == FocusBtn && msg.String() == "enter":
			m.Output.SetValue("")
			var res []rune
			decrypt := (m.Mode == Decrypt)

			var inputSource []rune
			if m.InputMode == File {
				if m.SelectedFile == "" {
					m.Output.SetValue("No file selected")
					return m, nil
				}
				content, err := os.ReadFile(m.SelectedFile)
				if err != nil {
					m.Output.SetValue("Error reading file: " + err.Error())
					return m, nil
				}
				inputSource = []rune(string(content))
			} else {
				inputSource = []rune(m.Input.Value())
			}

			switch m.Method {
			case Caesar:
				key, err := strconv.Atoi(m.Shift.Value())
				if err != nil {
					m.Output.SetValue("Shift should be positive number")
					return m, nil
				}
				if decrypt {
					key = -key
				}
				res = crypto.Caesar(inputSource, key)
			case Vigenere:
				if _, err := strconv.Atoi(m.Key.Value()); err == nil {
					m.Output.SetValue("Key cannot be a number")
					return m, nil
				}
				res = crypto.Vigenere(inputSource, []rune(m.Key.Value()), decrypt)
			}

			m.Output.SetValue(string(res))
			return m, nil
		case m.Focus == FocusSaveBtn && msg.String() == "enter":
			m.ShowSavePopup = true
			m.SaveFilename.SetValue("")
			m.SaveFilename.Focus()
			return m, nil
		}
	}

	var cmd tea.Cmd

	if m.Focus == FocusInput {
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.Focus == FocusFilepicker {
		m.Filepicker, cmd = m.Filepicker.Update(msg)
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

func (m model) View() string {
	layout := NewLayout(m.Width, m.Height)

	inputModeLabel := labelStyle.Render("INPUT MODE")
	inputModeValue := "File"
	if m.InputMode == ManualInput {
		inputModeValue = "Manual input"
	}
	if m.Focus == FocusInputMode {
		inputModeValue = selectedStyle.Render(inputModeValue) + " (Enter)"
	}

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

	label := "SAVE TO FILE"
	if m.Focus == FocusSaveBtn {
		label = saveBtnStyle.Render(label)
	}
	saveBtnLabel := label

	leftPanel := lipgloss.JoinVertical(
		lipgloss.Left,
		inputModeLabel,
		inputModeValue,
		"",
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
		saveBtnLabel,
	)

	leftPanel = panelStyle(LeftPanelWidth, layout.LeftPanelHeight).Render(lipgloss.PlaceVertical(layout.LeftPanelHeight, lipgloss.Top, leftPanel))

	inputPanelStyle := lipgloss.NewStyle().
		Height(layout.TopPanelHeight).
		Width(layout.RightPanelWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1)

	outputPanelStyle := lipgloss.NewStyle().
		Height(layout.BottomPanelHeight).
		Width(layout.RightPanelWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1)

	inputTitle := "INPUT"
	inputContent := m.Input.View()

	if m.InputMode == File {
		if m.FileConfirmed {
			inputTitle = "FILE SELECTED"

			style := selectedStyle
			if m.Focus != FocusFilepicker {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
			}

			selectedLine := style.Render("Selected: " + m.SelectedFile)
			hintLine := lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Italic(true).
				Render("Press Backspace to change file")
			inputContent = lipgloss.JoinVertical(lipgloss.Left, selectedLine, "", hintLine)
		} else {
			inputTitle = "FILE PICKER"
			inputContent = m.Filepicker.View()
		}
	}

	inputPanel := inputPanelStyle.Render(lipgloss.PlaceVertical(layout.TopPanelHeight, lipgloss.Top, labelStyle.Render(inputTitle)+"\n"+inputContent))

	outputView := m.Output.View()
	if m.Status != "" {
		outputView = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			Render(m.Status)
	}

	outputPanel := outputPanelStyle.Render(lipgloss.PlaceVertical(layout.BottomPanelHeight, lipgloss.Top, labelStyle.Render("OUTPUT")+"\n"+outputView))

	rightPanel := lipgloss.JoinVertical(0, inputPanel, outputPanel)

	baseView := lipgloss.JoinHorizontal(0, leftPanel, rightPanel)

	if m.ShowSavePopup {
		popupContent := lipgloss.JoinVertical(
			lipgloss.Left,
			labelStyle.Render("SAVE OUTPUT TO FILE"),
			"",
			"Filename:",
			m.SaveFilename.View(),
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Italic(true).
				Render("Enter to save  •  Esc to cancel"),
		)
		popup := popupStyle.Render(popupContent)
		return lipgloss.Place(
			m.Width, m.Height,
			lipgloss.Center, lipgloss.Center,
			popup,
			lipgloss.WithWhitespaceChars(" "),
			lipgloss.WithWhitespaceForeground(lipgloss.Color("237")),
		)
	}

	return baseView
}

func (m *model) skipInactive(dir int) {
	for {
		inactive := false
		if m.InputMode == ManualInput && m.Focus == FocusFilepicker {
			inactive = true
		}
		if m.InputMode == File && m.Focus == FocusInput {
			inactive = true
		}
		if m.Method == Caesar && m.Focus == FocusKey {
			inactive = true
		}
		if m.Method == Vigenere && m.Focus == FocusShift {
			inactive = true
		}

		if !inactive {
			break
		}
		m.Focus = (m.Focus + dir + FocusCount) % FocusCount
	}
}

func (m *model) updateFocus() {
	m.Input.Blur()
	m.Key.Blur()
	m.Shift.Blur()
	m.SaveFilename.Blur()

	switch m.Focus {
	case FocusInput:
		m.Input.Focus()
	case FocusKey:
		m.Key.Focus()
	case FocusShift:
		m.Shift.Focus()
	}

	if m.InputMode == File {
		m.Input.Blur()
	}

	switch m.Method {
	case Caesar:
		m.Key.SetValue("")
		m.Key.Blur()
	case Vigenere:
		m.Shift.SetValue("")
		m.Shift.Blur()
	}
}
