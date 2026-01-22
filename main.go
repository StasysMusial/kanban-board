package main

import (
	"fmt"
	"math/rand"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type Mode int
const(
	MODE_NORMAL Mode = iota
	MODE_EDIT
	MODE_TEXTEDIT
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

	mode      Mode
	cursor    int
	tags      []Tag
	boards    []Board
	editor    Editor
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
		mode: MODE_NORMAL,
		help: Help{},
		cursor: 0,
		can_paste: false,
	}

	InitKeyContexts()
	m.help.SetKeyContext(KEY_CONTEXT_BOARDS)

	// set up editor
	m.editor = NewEditor(&m)

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
	for i := range 4 {
		board := NewBoard(fmt.Sprintf("BOARD #%d", i+1), b_colors[i])
		m.boards = append(m.boards, board)
	}

	// create tasks
	for i := range len(m.boards) {
		tasks := []Task{}
		for j := range 3 {
			task := NewTask(&m, fmt.Sprintf("Task #%d", j+1))
			task.description = "Description"
			task.tags = j+i
			tasks = append(tasks, task)
		}
		m.boards[i].tasks = tasks
	}

	return m
}

func (m *model) SendTaskToEditor(task Task) {
	m.editor = m.editor.LoadTask(task)
}

func (m model) GetBoardCount() int {
	return len(m.boards)
}

func (m *model) Print(message string, color lip.Color) {
	m.message = lip.NewStyle().
		Foreground(color).
		Render(fmt.Sprintf(" → %s", message))
}

func (m *model) IncCursor() {
	if m.cursor < len(m.boards)-1 {
		m.cursor++
	}
	for i := range m.boards {
		selected := (i == m.cursor)
		m.boards[i].SetSelected(selected)
	}
}

func (m *model) DecrCursor() {
	if m.cursor > 0 {
		m.cursor--
	}
	for i := range m.boards {
		selected := (i == m.cursor)
		m.boards[i].SetSelected(selected)
	}
}

func (m *model) AddTask() {
	task := NewTask(m, "New Task")
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
	m.boards[m.cursor].AddTask(m.clipboard, board.cursor+offset)
	if below {
		m.boards[m.cursor].IncCursor()
	}
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

	m.boards[index].AddTask(task, 0)
}

func (m model) GetSelectedTask() Task {
	board := m.boards[m.cursor]
	if board.IsEmpty() {
		return Task{}
	}
	task := board.tasks[board.cursor]
	return task
}

func (m *model) SetSelectedTask(task Task) {
	board := m.boards[m.cursor]
	board.tasks[board.cursor] = task
	m.boards[m.cursor] = board
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
		boardsCount := m.GetBoardCount()
		boardsAreaWidth := m.width-EDITOR_WIDTH
		width := boardsAreaWidth / len(m.boards)
		width -= 2
		for i := range boardsCount {
			m.boards[i].width = width
			m.boards[i].SetHeight(m.height-7)
		}
	// process input
	case tea.KeyMsg:
		switch m.mode {
		case MODE_EDIT:
			switch msg.String() {
			case "enter":
				m.mode = MODE_NORMAL
				m.editor.name.Blur()
				task := m.GetSelectedTask()
				task.name = m.editor.name.Value()
				m.SetSelectedTask(task)
			}
			// edit mode mappings go here
		case MODE_NORMAL:
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
			case "enter":
				m.mode = MODE_EDIT
				m.editor.name.Focus()
				m.editor.name.CursorEnd()
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
	}

	// update boards
	for i := range len(m.boards) {
		var cmd tea.Cmd
		m.boards[i], cmd = m.boards[i].Update(msg, m)
		m.boards[i].selected = (i == m.cursor)
		cmds = append(cmds, cmd)
	}

	// update editor
	var editorCmd tea.Cmd
	m.editor, editorCmd = m.editor.Update(msg, m)
	cmds = append(cmds, editorCmd)

	if m.mode == MODE_NORMAL {
		m.SendTaskToEditor(m.GetSelectedTask())
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
		boards = append(boards, board.View(m))
	}
	result := lip.JoinHorizontal(
		lip.Top,
		boards...
	)
	result = lip.JoinHorizontal(
		lip.Top,
		result,
		m.editor.View(),
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
