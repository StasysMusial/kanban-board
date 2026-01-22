package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

const EDITOR_WIDTH = 40

type Editor struct {
	width  int
	height int

	task   Task
	name   textinput.Model

	mptr   *model
}

func NewEditor(mptr *model) Editor {
	nameInput := textinput.New()
	nameInput.Placeholder = ""
	nameInput.Width = EDITOR_WIDTH-8
	nameInput.Prompt = ""
	nameInput.CharLimit = EDITOR_WIDTH-8
	editor := Editor{
		name: nameInput,
		mptr: mptr,
	}
	return editor
}

func (e Editor) LoadTask(task Task) Editor {
	e.task = task
	e.name.SetValue(task.name)
	return e
}

func (e Editor) Init() tea.Cmd {
	return nil
}

func (e Editor) Update(msg tea.Msg, m model) (Editor, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.mode != MODE_EDIT {
			break
		}
		var nameCmd tea.Cmd
		e.name, nameCmd = e.name.Update(msg)
		cmds = append(cmds, nameCmd)
	}
	return e, tea.Batch(cmds...)
}

func (e Editor) View() string {
	name := lip.JoinVertical(
		lip.Left,
		"Task:",
		e.name.View(),
	)
	result := editorStyle.containerStyle.Render(name)
	return result
}
