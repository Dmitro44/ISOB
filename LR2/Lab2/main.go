package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(NewModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}
