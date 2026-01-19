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

func (m *model) IncCursor() {
	if m.cursor < len(m.boards)-1 {
		m.cursor++
	} else {
		// m.cursor = 0
	}
	for i := range m.boards {
		selected := (i == m.cursor)
		m.boards[i].SetSelected(selected)
	}
}

func (m *model) DecrCursor() {
	if m.cursor > 0 {
		m.cursor--
	} else {
		// m.cursor = len(m.boards)-1
	}
	for i := range m.boards {
		selected := (i == m.cursor)
		m.boards[i].SetSelected(selected)
	}
}

func initialModel() tea.Model {
	var m model = model{
		tags: []Tag{
			NewTag("󰫢", lip.Color("#ff4cc4")),
			NewTag("󰅩", lip.Color("#89d789")),
			NewTag("󰃣", lip.Color("#f5d33d")),
			NewTag("", lip.Color("#5c84d6")),
		},
		cursor: 0,
	}
	m.boards = []Board{}
	b_colors := []lip.Color{
		lip.Color("1"),
		lip.Color("2"),
		lip.Color("3"),
		lip.Color("4"),
	}
	for i := range 4 {
		board := NewBoard(&m, fmt.Sprintf("BOARD #%d", i+1), b_colors[i])
		m.boards = append(m.boards, board)
	}
	for i := range len(m.boards) {
		tasks := []Task{}
		for j := range 14+i {
			task := NewTask(&m, &m.boards[i], fmt.Sprintf("Task #%d", j+1))
			task.description = "Description"
			task.tags = j+i
			tasks = append(tasks, task)
		}
		m.boards[i].tasks = tasks
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
		width := m.width / len(m.boards)
		for i := range len(m.boards) {
			m.boards[i].width = width-2
			m.boards[i].SetHeight(m.height-5)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h", "left":
			m.DecrCursor()
		case "l", "right":
			m.IncCursor()
		}
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
	var boards []string
	for _, board := range m.boards {
		boards = append(boards, board.View())
	}
	result := lip.JoinHorizontal(
		lip.Top,
		boards...
	)
	result = lip.JoinVertical(
		lip.Left,
		result,
		uiStyle.helpStyle.Width(m.width).AlignHorizontal(lip.Center).Render("q - quit   h/j/k/l - navigate"),
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
