package main

import (
	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type Board struct {
	width    int
	height   int

	title    string
	color    lip.Color

	selected bool
	scroll   int
	cursor   int
	tasks    []Task

	mptr     *model
}

func NewBoard(mptr *model, title string, color lip.Color) Board {
	return Board{
		title:    title,
		color:    color,

		selected: false,
		scroll:   0,
		cursor:   0,
		tasks:    []Task{},
		mptr:     mptr,
	}
}

func (b Board) Init() tea.Cmd {
	return nil
}

func (b Board) Update(msg tea.Msg) (Board, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.height = msg.Height
	}

	// Update tasks
	for i := range len(b.tasks) {
		var cmd tea.Cmd
		b.tasks[i], cmd = b.tasks[i].Update(msg)
		cmds = append(cmds, cmd)
		b.tasks[i].selected = (i == b.cursor)
	}

	return b, tea.Batch(cmds...)
}

func (b Board) View() string {
	var tasks []string
	for _, task := range b.tasks {
		tasks = append(tasks, task.View())
	}

	result := lip.JoinVertical(
		lip.Left,
		tasks...
	)

	result = lip.Place(
		b.width,
		b.height,
		lip.Center,
		lip.Center,
		result,
	)

	return result
}
