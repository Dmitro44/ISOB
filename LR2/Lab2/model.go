package main

import (
	"go-cipher/crypto"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("237"))

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

	borderColor = lipgloss.Color("99")

	saveBtnStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("86")).
			Foreground(lipgloss.Color("0")).
			Padding(0, 2).
			Bold(true)

	popupStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("86")).
			Padding(1, 2).
			Width(40)
)

func panelStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1).
		Width(width).
		Height(height)
}

type (
	Method    int
	Mode      int
	InputMode int
)

const (
	ManualInput InputMode = iota
	File
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
	FocusInputMode int = iota
	FocusFilepicker
	FocusMethod
	FocusMode
	FocusKey
	FocusShift
	FocusInput
	FocusBtn
	FocusSaveBtn
	FocusCount
)

type model struct {
	Filepicker    filepicker.Model
	SelectedFile  string
	FileConfirmed bool
	InputMode     InputMode
	Method        Method
	Mode          Mode
	Focus         int
	Width         int
	Height        int
	Key           textinput.Model
	Shift         textinput.Model
	Input         textarea.Model
	Output        textarea.Model
	SaveFilename  textinput.Model
	ShowSavePopup bool
	keymap        keymap
}

type keymap = struct {
	next, prev, quit, copy, clearInput key.Binding
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

func newFilepicker() filepicker.Model {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".txt", ".md", ".go"}
	dir, _ := os.Getwd()
	fp.CurrentDirectory = dir
	fp.FileAllowed = true
	fp.DirAllowed = false
	return fp
}

func NewModel() model {
	m := model{
		Filepicker:   newFilepicker(),
		Method:       Caesar,
		Mode:         Encrypt,
		Focus:        FocusInput,
		Input:        newTextarea("Type text you want to cipher..."),
		Output:       newOutputTextarea("Result will appear here..."),
		SaveFilename: newTextinput("output.txt"),
		Key:          newTextinput("Enter key..."),
		Shift:        newTextinput("3"),
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
			copy: key.NewBinding(
				key.WithKeys("ctrl+d"),
			),
			clearInput: key.NewBinding(
				key.WithKeys("alt+d"),
			),
		},
	}

	return m
}

func (m model) Init() tea.Cmd {
	return m.Filepicker.Init()
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
		rightContentW := msg.Width - leftWidth - 8

		// textarea overhead inside panel: 1 (label) + 2 (textarea border) = 3
		m.Input.SetWidth(rightContentW - 2) // -2 for textarea border left+right
		m.Input.SetHeight(topContentH - 3)
		m.Output.SetWidth(rightContentW - 2)
		m.Output.SetHeight(botContentH - 3)
		// Set filepicker dimensions
		m.Filepicker.SetHeight(topContentH - 2)
		m.SaveFilename.Width = 34
	case tea.KeyMsg:
		// ── Popup intercepts ALL keys when visible ──────────────────────────
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
				err := os.WriteFile(filename, []byte(m.Output.Value()), 0644)
				if err != nil {
					m.Output.SetValue("Error saving file: " + err.Error())
				} else {
					m.Output.SetValue("Saved to: " + filename)
				}
				m.ShowSavePopup = false
				m.SaveFilename.Blur()
				return m, nil
			default:
				var cmd tea.Cmd
				m.SaveFilename, cmd = m.SaveFilename.Update(msg)
				return m, cmd
			}
		}
		// ── End popup handler ───────────────────────────────────────────────

		// Handle global keys first (Quit, Tab)
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

		// Handle Filepicker specific keys
		if m.InputMode == File && m.Focus == FocusFilepicker && !m.FileConfirmed {
			var cmd tea.Cmd
			m.Filepicker, cmd = m.Filepicker.Update(msg)
			if didSelect, path := m.Filepicker.DidSelectFile(msg); didSelect {
				m.SelectedFile = path
				m.FileConfirmed = true
			}
			return m, cmd
		}

		if m.InputMode == File && m.Focus == FocusFilepicker && m.FileConfirmed {
			if msg.String() == "backspace" {
				m.FileConfirmed = false
				m.SelectedFile = ""
				return m, m.Filepicker.Init()
			}
			return m, nil
		}

		// Handle other focused elements
		switch {
		case key.Matches(msg, m.keymap.copy):
			m.Input.SetValue(m.Output.Value())
			m.Output.SetValue("")
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
			// if currently focused field became inactive, skip it
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

func (m model) View() string {
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

	leftWidth := 25
	rightWidth := m.Width - leftWidth - 8

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
		BorderForeground(borderColor).
		Padding(1)

	outputPanelStyle := lipgloss.NewStyle().
		Height(botH).
		Width(rightWidth).
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

	inputPanel := inputPanelStyle.Render(labelStyle.Render(inputTitle) + "\n" + inputContent)

	outputPanel := outputPanelStyle.Render(labelStyle.Render("OUTPUT") + "\n" + m.Output.View())

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
