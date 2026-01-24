package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type Mode int
const (
	MODE_NORMAL Mode = iota
	MODE_EDIT
)

const MIN_WIN_WIDTH = 96

type model struct {
	width     int
	height    int

	mode      Mode
	cursor    int
	boards    []Board
	editor    Editor
	help      Help

	clipboard Task
	message   string
	can_paste bool
}

func initialModel() tea.Model {
	InitStyles()
	InitKeyContexts()

	var m model = model{
		mode: MODE_NORMAL,
		help: Help{},
		cursor: 0,
		can_paste: false,
	}

	m.help.SetKeyContext(KEY_CONTEXT_BOARDS)

	// set up editor
	m.editor = NewEditor()

	// TODO: this needs to be using the
	// per project config instead of the default config
	conf := GetConfig(GetDefaultConfigPath())
	m.LoadConfig(conf)

	return m
}

func (m *model) LoadConfig(c config) {
	for _, t := range c.Tags {
		NewTag(t.Icon, lip.Color(t.Color))
	}
	for _, b := range c.Boards {
		board := NewBoard(b.Name, lip.Color(b.Color))
		m.boards = append(m.boards, board)
	}
}

func (m *model) SendTaskToEditor(task Task) {
	m.editor = m.editor.LoadTask(task)
}

func (m model) GetBoardCount() int {
	return len(m.boards)
}

func (m model) HasBoards() bool {
	return m.GetBoardCount() > 0
}

func (m *model) Print(message string, color lip.Color) {
	m.message = lip.NewStyle().
		Foreground(color).
		Render(fmt.Sprintf("→ %s", message))
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
	task := NewTask(m, "")
	m.boards[m.cursor].AddTask(task, m.boards[m.cursor].cursor+1)
	m.boards[m.cursor].IncCursor()
	m.Print(fmt.Sprintf("Added task to [%s]", m.boards[m.cursor].title), msgColorInfo)
}

func (m *model) CopyTask() {
	board := m.boards[m.cursor]
	m.clipboard = board.tasks[board.cursor]
	m.can_paste = true
	m.Print(fmt.Sprintf("Yanked [%s]", m.clipboard.name), msgColorInfo)
}

func (m *model) RemoveTask() {
	board := m.boards[m.cursor]
	m.boards[m.cursor].RemoveTask(board.cursor)
	m.Print(fmt.Sprintf("Cut [%s]", m.clipboard.name), msgColorInfo)
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
	m.Print(fmt.Sprintf("Pasted [%s] to [%s]", m.clipboard.name, board.title), msgColorInfo)
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
	if len(m.boards) < 1 {
		return Task{}
	}
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

func (m *model) EnterModeEdit() {
	m.mode = MODE_EDIT
	m.editor.name.Focus()
	m.editor.name.CursorEnd()
	m.editor.mode = EDIT_MODE_NAME
}

func (m *model) EnterModeNormal() {
	m.mode = MODE_NORMAL
	m.editor.name.Blur()
	m.editor.desc.Blur()
	m.help.context_data = keyContextBoards
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	deferredEdit := false

	switch msg := msg.(type) {
	// process window size
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.HasBoards() {
			break
		}
		boardsCount := m.GetBoardCount()
		boardsAreaWidth := m.width-EDITOR_WIDTH
		width := boardsAreaWidth / len(m.boards)
		width -= 1
		for i := range boardsCount {
			m.boards[i].width = width
			m.boards[i].SetHeight(m.height-7)
		}
	// process input
	case tea.KeyMsg:
		switch m.mode {
		case MODE_EDIT:
			if m.editor.mode != EDIT_MODE_NAME {
				break
			}
			if !m.HasBoards() {
				break
			}
			switch msg.String() {
			case "enter":
				m.EnterModeNormal()
				task := m.GetSelectedTask()
				task.name = m.editor.name.Value()
				task.description = m.editor.desc.Value()
				task.tags = m.editor.tags
				m.SetSelectedTask(task)
			case "esc":
				m.EnterModeNormal()
			}
			// edit mode mappings go here
		case MODE_NORMAL:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			}

			if !m.HasBoards() {
				break
			}

			switch msg.String() {
			case "a":
				m.AddTask()
				deferredEdit = true
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
				case "enter":
					m.EnterModeEdit()
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

	if deferredEdit {
		m.EnterModeEdit()
	}

	if m.mode == MODE_EDIT {
		switch m.editor.mode {
		case EDIT_MODE_NAME:
			m.help.context_data = keyContextTask
		case EDIT_MODE_DESC:
			m.help.context_data = keyContextTaskDesc
		}
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
	result := ""
	if m.HasBoards() {
		result = lip.JoinHorizontal(
			lip.Left,
			boards...
		)
	} else {
		result = lip.Place(
			m.width - EDITOR_WIDTH,
			m.height - 2,
			lip.Center,
			lip.Center,
			"NO BOARDS\nconfigure boards in .kanban/config.toml\npress q to quit",
		)
	}
	result = lip.JoinVertical(
		lip.Left,
		result,
		m.message,
		m.help.View(m, lip.Width(result)),
	)
	result = lip.JoinHorizontal(
		lip.Top,
		result,
		m.editor.View(m),
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
