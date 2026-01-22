package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type EditMode int
const(
	EDIT_MODE_NAME EditMode = iota
	EDIT_MODE_DESC
)

const EDITOR_WIDTH = 40

type Editor struct {
	width  int
	height int

	task   Task
	name   textinput.Model
	desc   textarea.Model
	mode   EditMode
}

func NewEditor() Editor {
	nameInput := textinput.New()
	nameInput.Placeholder = ""
	nameInput.Width = EDITOR_WIDTH-8
	nameInput.Prompt = ""
	nameInput.CharLimit = EDITOR_WIDTH-8

	descInput := textarea.New()
	descInput.Prompt = ""
	descInput.SetWidth(EDITOR_WIDTH-4)
	descInput.Placeholder = ""
	descInput.SetHeight(10)
	descInput.ShowLineNumbers = false
	editor := Editor{
		name: nameInput,
		desc: descInput,
		mode: EDIT_MODE_NAME,
	}
	return editor
}

func (e Editor) LoadTask(task Task) Editor {
	e.task = task
	e.name.SetValue(task.name)
	e.desc.SetValue(task.description)
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

		switch msg.String() {
		case "tab":
			switch e.mode {
			case EDIT_MODE_NAME:
				e.mode = EDIT_MODE_DESC
				e.name.Blur()
				e.desc.Focus()
			case EDIT_MODE_DESC:
				e.mode = EDIT_MODE_NAME
				e.desc.Blur()
				e.name.Focus()
			}
		}

		var nameCmd tea.Cmd
		e.name, nameCmd = e.name.Update(msg)
		cmds = append(cmds, nameCmd)

		var descCmd tea.Cmd
		e.desc, descCmd = e.desc.Update(msg)
		cmds = append(cmds, descCmd)
	}
	return e, tea.Batch(cmds...)
}

func (e Editor) View(m model) string {
	var style EditorStyle
	if m.mode == MODE_NORMAL {
		style = editorStyle
	} else {
		switch e.mode {
		case EDIT_MODE_NAME:
			style = editorStyleName
		case EDIT_MODE_DESC:
			style = editorStyleDesc
		}
	}
	name := e.name.View()
	desc := e.desc.View()
	result := lip.JoinVertical(
		lip.Left,
		style.nameLabelStyle.Render(" Task "),
		"",
		name,
		"",
		style.descLabelStyle.Render(" Description "),
		"",
		desc,
	)
	result = editorStyle.containerStyle.Render(result)
	return result
}
