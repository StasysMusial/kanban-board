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
	lists []List
}

func initialModel() tea.Model {
	lists := []List{}
	for j := range 3 {
		var content []string = []string{}
		for i := range 64 {
			content = append(
				content,
				fmt.Sprintf("test line %s", strconv.Itoa(i)),
			)
		}
		list := NewList(30, 10, content)
		lists = append(lists, list)
	}
	return model{
		lists: lists,
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

	for _, list := range m.lists {
		list, cmd = list.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}
