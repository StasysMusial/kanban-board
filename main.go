package main

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

}

type model struct {
	list List
}

func initialModel() tea.Model {
	var content []string = []string{}
	for i := range 64 {
		content = append(
			content,
			fmt.Sprintf("test line %s", strconv.Itoa(i)),
		)
	}
	return model{
		list: NewList(30, 10, content),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}
