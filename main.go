package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

func main() {
	InitStyles()
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	width  int
	height int

	cursor int
	tags   []Tag
	boards []Board
}

func initialModel() tea.Model {
	var m model = model{
		tags:     []Tag{
			NewTag("󰫢", lip.Color("#ff4cc4")),
			NewTag("󰅩", lip.Color("#89d789")),
			NewTag("󰃣", lip.Color("#f5d33d")),
			NewTag("", lip.Color("#5c84d6")),
		},
		cursor: 0,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
	}

	for i := range len(m.boards) {
		var cmd tea.Cmd
		m.boards[i], cmd = m.boards[i].Update(msg)
		m.boards[i].selected = (i == m.cursor)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return ""
}
