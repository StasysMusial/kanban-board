package main

import (
	"fmt"
	"os"
	"strconv"

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
	tasks  []Task
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
	for i := range 7 {
		task := NewTask(&m, fmt.Sprintf("Task #%s", strconv.Itoa(i+1)))
		task.description = "Description"
		task.tags = i + 7
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
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "p":
			fmt.Println(m.tasks[m.cursor].GetTags())
		}
	}

	for i := range len(m.tasks) {
		m.tasks[i], cmd = m.tasks[i].Update(msg)
		m.tasks[i].selected = (i == m.cursor)
	}

	return m, cmd
}

func (m model) View() string {
	var tasks []string
	for _, task := range m.tasks {
		tasks = append(tasks, task.View())
	}

	result := lip.JoinVertical(
		lip.Left,
		tasks...
	)

	result = lip.Place(
		m.width,
		m.height,
		lip.Center,
		lip.Center,
		result,
	)

	return result
}
