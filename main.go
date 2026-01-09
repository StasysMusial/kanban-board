package main

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

}

type model struct {
	width  int
	height int
	lists  []List
}

func initialModel() tea.Model {
	lists := []List{}
	for range 3 {
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	for i := range m.lists {
		m.lists[i], cmd = m.lists[i].Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	lists := []string{}
	for _, list := range m.lists {
		lists = append(lists, list.View())
	}
	s := lip.JoinHorizontal(lip.Center, lists...)

	return lip.Place(
		m.width,
		m.height,
		lip.Center,
		lip.Center,
		s,
	)
}
