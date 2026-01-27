package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Println(fmt.Sprintf("kanban-board v%s", CURRENT_VERSION))
			return
		}
	}

	InitVersion()
	InitIO()
	GenerateDefaultConfig()
	if !GenerateProjectConfig() {
		return
	}
	// set up main model
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
