package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	switch len(os.Args) {
	case 1:
		break
	case 2:
		arg := os.Args[1]
		switch arg {
		case "-v", "--version":
			fmt.Printf("kanban-board v%s\n", CURRENT_VERSION)
			return
		default:
			fmt.Printf("Error: Unrecognized argument '%s'\n", arg)
			return
		}
	default:
		fmt.Println("Error: Too many arguments")
		return
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
