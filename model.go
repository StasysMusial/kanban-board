package main

import (
	"fmt"
	"strings"

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
	columns   []Column
	editor    Editor
	help      Help

	clipboard Task
	message   string
	title     string
	color     lip.Color
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

	m.help.SetKeyContext(KEY_CONTEXT_COLUMNS)

	// set up editor
	m.editor = NewEditor()

	// TODO: this needs to be using the
	// per project config instead of the default config
	conf := GetConfig(GetCwdConfigPath())
	m.LoadConfig(conf)

	data := ReadData()
	m.LoadState(data)

	return m
}

func (m *model) LoadConfig(c config) {
	for _, t := range c.Tags {
		NewTag(t.Icon, lip.Color(t.Color))
	}
	for _, b := range c.Columns {
		board := NewColumn(b.Name, b.Icon, lip.Color(b.Color))
		m.columns = append(m.columns, board)
	}
	title := c.Title
	if title == "DEFAULT" {
		cwdParts := strings.Split(cwd, pathSeparator)
		title = cwdParts[len(cwdParts)-1]
	}
	m.title = title
	color := c.Color
	if color == "DEFAULT" {
		color = "7"
	}
	m.color = lip.Color(color)
}

func (m *model) LoadState(d modelSaveData) {
	for i, bData := range d.Columns {
		m.columns[i].tasks = []Task{}
		for _, t := range bData.Tasks {
			task := Task{
				name: t.Name,
				description: t.Desc,
				tags: t.Tags,
			}
			if i < len(m.columns) {
				m.columns[i].tasks = append(m.columns[i].tasks, task)
			}
		}
	}
}

func (m *model) SendTaskToEditor(task Task) {
	m.editor = m.editor.LoadTask(task)
}

func (m model) GetColumnCount() int {
	return len(m.columns)
}

func (m model) HasColumns() bool {
	return m.GetColumnCount() > 0
}

func (m *model) Print(message string, color lip.Color) {
	m.message = lip.NewStyle().
		Foreground(color).
		Render(fmt.Sprintf("→ %s", message))
}

func (m *model) IncCursor() {
	if m.cursor < len(m.columns)-1 {
		m.cursor++
	}
	for i := range m.columns {
		selected := (i == m.cursor)
		m.columns[i].SetSelected(selected)
	}
}

func (m *model) DecrCursor() {
	if m.cursor > 0 {
		m.cursor--
	}
	for i := range m.columns {
		selected := (i == m.cursor)
		m.columns[i].SetSelected(selected)
	}
}

func (m *model) AddTask() {
	task := NewTask(m, "")
	m.columns[m.cursor].AddTask(task, m.columns[m.cursor].cursor+1)
	m.columns[m.cursor].IncCursor()
	m.Print(fmt.Sprintf("Added task to [%s]", m.columns[m.cursor].title), msgColorInfo)
}

func (m *model) CopyTask() {
	board := m.columns[m.cursor]
	m.clipboard = board.tasks[board.cursor]
	m.can_paste = true
	m.Print(fmt.Sprintf("Yanked [%s]", m.clipboard.name), msgColorInfo)
}

func (m *model) RemoveTask() {
	board := m.columns[m.cursor]
	m.columns[m.cursor].RemoveTask(board.cursor)
}

func (m *model) CutTask() {
	m.CopyTask()
	m.RemoveTask()
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
	board := m.columns[m.cursor]
	m.columns[m.cursor].AddTask(m.clipboard, board.cursor+offset)
	if below {
		m.columns[m.cursor].IncCursor()
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
	if index >= len(m.columns) {
		return
	}
	board := m.columns[m.cursor]
	task := board.tasks[board.cursor]
	m.columns[m.cursor].RemoveTask(board.cursor)

	m.columns[index].AddTask(task, 0)
	m.Print(fmt.Sprintf("Moved [%s] to [%s]", task.name, m.columns[index].title), msgColorInfo)
}

func (m *model) MoveTaskLeft() {
	index := m.cursor-1
	if index < 0 {
		return
	}
	board := m.columns[m.cursor]
	task := board.tasks[board.cursor]
	m.columns[m.cursor].RemoveTask(board.cursor)

	m.columns[index].AddTask(task, 0)
	m.Print(fmt.Sprintf("Moved [%s] to [%s]", task.name, m.columns[index].title), msgColorInfo)
}

func (m model) GetSelectedTask() Task {
	if len(m.columns) < 1 {
		return Task{}
	}
	column := m.columns[m.cursor]
	column.ClampCursor()
	if column.IsEmpty() {
		return Task{}
	}
	task := column.tasks[column.cursor]
	return task
}

func (m *model) SetSelectedTask(task Task) {
	column := m.columns[m.cursor]
	column.ClampCursor()
	column.tasks[column.cursor] = task
	m.columns[m.cursor] = column
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
	m.help.context_data = keyContextBoard
}

func (m model) ViewTitle() string {
	return projectTitleStyle.
		MaxHeight(1).
		Background(m.color).
		Render(fmt.Sprintf(" %s ", m.title))
}

func (m model) ViewMessage(width int) string {
	return lip.NewStyle().MaxWidth(width).Render(m.message)
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
		if !m.HasColumns() {
			break
		}
		boardsCount := m.GetColumnCount()
		boardsAreaWidth := m.width-EDITOR_WIDTH
		width := boardsAreaWidth / len(m.columns)
		width -= 1
		for i := range boardsCount {
			m.columns[i].width = width
			m.columns[i].SetHeight(m.height-7)
		}
	// process input
	case tea.KeyMsg:
		switch m.mode {
		case MODE_EDIT:
			if m.editor.mode != EDIT_MODE_NAME {
				break
			}
			if !m.HasColumns() {
				break
			}
			switch msg.String() {
			case "enter":
				m.EnterModeNormal()
				task := m.GetSelectedTask()
				task.name = m.editor.name.Value()
				task.name = strings.TrimSpace(task.name)
				task.description = m.editor.desc.Value()
				task.description = strings.TrimSpace(task.description)
				task.tags = m.editor.tags
				m.SetSelectedTask(task)
			case "esc":
				m.EnterModeNormal()
				if m.editor.task.IsEmpty() {
					m.RemoveTask()
					m.message = ""
				} else {
					m.SetSelectedTask(m.editor.task)
				}
				Undo(&m, false)
			}
			// edit mode mappings go here
		case MODE_NORMAL:
			switch msg.String() {
			case "q":
				data := ModelToJSON(m)
				WriteData(data)
				return m, tea.Quit
			case "esc":
				m.message = ""
			case "u":
				Undo(&m, true)
				for i := range m.columns {
					m.columns[i].ClampCursor()
				}
			case "ctrl+r":
				Redo(&m)
				for i := range m.columns {
					m.columns[i].ClampCursor()
				}
			}

			if !m.HasColumns() {
				break
			}

			switch msg.String() {
			case "a":
				AddUndoPoint(m, true)
				m.AddTask()
				deferredEdit = true
			case "p":
				AddUndoPoint(m, true)
				m.PasteTaskBelow()
			case "P":
				AddUndoPoint(m, true)
				m.PasteTaskAbove()
			case "h":
				m.DecrCursor()
			case "l":
				m.IncCursor()
			}
			// only do this stuff if selected board has tasks inside it
			if !m.columns[m.cursor].IsEmpty() {
				switch msg.String() {
				case "H":
					AddUndoPoint(m, true)
					m.MoveTaskLeft()
					m.DecrCursor()
				case "L":
					AddUndoPoint(m, true)
					m.MoveTaskRight()
					m.IncCursor()
				case "x":
					AddUndoPoint(m, true)
					m.CutTask()
				case "y":
					m.CopyTask()
				case "enter":
					AddUndoPoint(m, true)
					m.EnterModeEdit()
				}
			}
		}
	}

	// update boards
	for i := range len(m.columns) {
		var cmd tea.Cmd
		m.columns[i], cmd = m.columns[i].Update(msg, m)
		m.columns[i].selected = (i == m.cursor)
		cmds = append(cmds, cmd)
	}

	// update editor
	var editorCmd tea.Cmd
	m.editor, editorCmd = m.editor.Update(msg, m)
	cmds = append(cmds, editorCmd)

	if m.mode == MODE_NORMAL {
		m.SendTaskToEditor(m.GetSelectedTask())
	} else {
		task := m.GetSelectedTask()
		task.name = m.editor.name.Value()
		task.description = m.editor.desc.Value()
		task.tags = m.editor.tags
		m.SetSelectedTask(task)
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
	for _, board := range m.columns {
		boards = append(boards, board.View(m))
	}
	result := ""
	if m.HasColumns() {
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
		m.ViewMessage(lip.Width(result)),
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
