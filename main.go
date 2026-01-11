package main

import (
	"container/list"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
	// lip "github.com/charmbracelet/lipgloss"
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
	width    int
	height   int
	tags     list.List
	tasks    []Task
	selected int
}

func initialModel() tea.Model {
	tags := list.List{}
	tags.Init()
	var m model = model{
		tags:     tags,
		selected: 0,
	}
	for range 3 {
		task := NewTask(&m, "Name")
		task.description = "Description"
		m.tasks = append(m.tasks, task)
	}
	return m
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
		case "j", "down":
			if m.selected < len(m.tasks)-1 {
				m.selected++
			}
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
		}
	}

	for i := range len(m.tasks) {
		m.tasks[i], cmd = m.tasks[i].Update(msg)
		m.tasks[i].selected = (i == m.selected)
	}

	return m, cmd
}

func (m model) View() string {
	var tasks []string
	for _, task := range m.tasks {
		tasks = append(tasks, task.View())
	}
	return lip.JoinVertical(
		lip.Left,
		tasks...
	)
}
