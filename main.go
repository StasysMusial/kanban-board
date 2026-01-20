package main

import (
	"fmt"
	"math/rand"
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
	width     int
	height    int

	cursor    int
	tags      []Tag
	boards    []Board
	help      Help

	clipboard Task
	message   string
	can_paste bool
}

func initialModel() tea.Model {
	var m model = model{
		tags: []Tag{
			NewTag("󰃣", lip.Color("#f5d33d")),
			NewTag("", lip.Color("#5c84d6")),
			NewTag("󰅩", lip.Color("#89d789")),
			NewTag("󰫢", lip.Color("#ff4cc4")),
		},
		help: Help{},
		cursor: 0,
		can_paste: false,
	}

	InitKeyContexts()
	m.help.SetKeyContext(KEY_CONTEXT_BOARDS)

	// prepare board colors
	b_colors := []lip.Color{
		lip.Color("1"),
		lip.Color("2"),
		lip.Color("3"),
		lip.Color("4"),
		lip.Color("5"),
	}

	// set up boards
	m.boards = []Board{}
	for i := range 5 {
		board := NewBoard(&m, fmt.Sprintf("BOARD #%d", i+1), b_colors[i])
		m.boards = append(m.boards, board)
	}

	// create tasks
	for i := range len(m.boards) {
		tasks := []Task{}
		for j := range 3 {
			task := NewTask(&m, &m.boards[i], fmt.Sprintf("Task #%d", j+1))
			task.description = "Description"
			task.tags = j+i
			tasks = append(tasks, task)
		}
		m.boards[i].tasks = tasks
	}

	return m
}

func (m *model) Print(message string, color lip.Color) {
	m.message = lip.NewStyle().
		Foreground(color).
		Render(fmt.Sprintf("→ %s", message))
}

func (m *model) IncCursor() {
	if m.cursor < len(m.boards)-1 {
		m.cursor++
	// } else {
	// 	m.cursor = 0
	}
	for i := range m.boards {
		selected := (i == m.cursor)
		m.boards[i].SetSelected(selected)
	}
}

func (m *model) DecrCursor() {
	if m.cursor > 0 {
		m.cursor--
	// } else {
	// 	m.cursor = len(m.boards)-1
	}
	for i := range m.boards {
		selected := (i == m.cursor)
		m.boards[i].SetSelected(selected)
	}
}

func (m *model) AddTask() {
	task := NewTask(m, &m.boards[m.cursor], "New Task")
	m.boards[m.cursor].AddTask(task, m.boards[m.cursor].cursor+1)
	m.boards[m.cursor].IncCursor()
	m.boards[m.cursor].tasks[m.boards[m.cursor].cursor].tags = rand.Int() % 16
	m.Print(fmt.Sprintf("Added task to [%s]", m.boards[m.cursor].title), msgInfoColor)
}

func (m *model) CopyTask() {
	board := m.boards[m.cursor]
	m.clipboard = board.tasks[board.cursor]
	m.can_paste = true
	m.Print(fmt.Sprintf("Yanked [%s]", m.clipboard.name), msgInfoColor)
}

func (m *model) RemoveTask() {
	board := m.boards[m.cursor]
	m.boards[m.cursor].RemoveTask(board.cursor)
	m.Print(fmt.Sprintf("Cut [%s]", m.clipboard.name), msgInfoColor)
}

func (m *model) PasteTask(below bool) {
	if !m.can_paste {
		return
	}
	offset := 0
	if below {
		offset = 1
	}
	board := m.boards[m.cursor]
	m.clipboard.bptr = &m.boards[m.cursor]
	m.boards[m.cursor].AddTask(m.clipboard, board.cursor+offset)
	m.Print(fmt.Sprintf("Pasted [%s] to [%s]", m.clipboard.name, board.title), msgInfoColor)
}

func (m *model) PasteTaskBelow() {
	m.PasteTask(true)
}

func (m *model) PasteTaskAbove() {
	m.PasteTask(false)
}

func (m *model) MoveTaskRight() {
	index := m.cursor+1
	if index >= len(m.boards) {
		// index = 0
		return
	}
	board := m.boards[m.cursor]
	task := board.tasks[board.cursor]
	m.boards[m.cursor].RemoveTask(board.cursor)

	task.bptr = &m.boards[index]
	m.boards[index].AddTask(task, 0)
}

func (m *model) MoveTaskLeft() {
	index := m.cursor-1
	if index < 0 {
		// index = len(m.boards)-1
		return
	}
	board := m.boards[m.cursor]
	task := board.tasks[board.cursor]
	m.boards[m.cursor].RemoveTask(board.cursor)

	task.bptr = &m.boards[index]
	m.boards[index].AddTask(task, 0)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	// process window size
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		width := m.width / len(m.boards)
		for i := range len(m.boards) {
			m.boards[i].width = width-2
			m.boards[i].SetHeight(m.height-7)
		}
	// process input
	case tea.KeyMsg:
		// m.Print("", lip.Color("0"))
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "a":
			m.AddTask()
		case "p":
			m.PasteTaskBelow()
		case "P":
			m.PasteTaskAbove()
		case "h":
			m.DecrCursor()
		case "l":
			m.IncCursor()
		}
		// only do this stuff if selected board has tasks inside it
		if !m.boards[m.cursor].IsEmpty() {
			switch msg.String() {
			case "H":
				m.MoveTaskLeft()
				m.DecrCursor()
			case "L":
				m.MoveTaskRight()
				m.IncCursor()
			case "x":
				m.CopyTask()
				m.RemoveTask()
			case "y":
				m.CopyTask()
			}
		}
	}

	// update boards
	for i := range len(m.boards) {
		var cmd tea.Cmd
		m.boards[i], cmd = m.boards[i].Update(msg)
		m.boards[i].selected = (i == m.cursor)
		cmds = append(cmds, cmd)
	}

	// update help section
	var helpCmd tea.Cmd
	m.help, helpCmd = m.help.Update(msg)
	cmds = append(cmds, helpCmd)

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
		m.message,
		m.help.View(),
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
