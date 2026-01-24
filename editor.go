package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
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
	tags   int
	mode   EditMode
}

func NewEditor() Editor {
	nameInput := textinput.New()
	nameInput.Placeholder = ""
	nameInput.Width = EDITOR_WIDTH-7
	nameInput.Prompt = ""
	nameInput.CharLimit = EDITOR_WIDTH-7

	descInput := textarea.New()
	descInput.Prompt = ""
	descInput.SetWidth(EDITOR_WIDTH-3)
	descInput.Placeholder = ""
	descInput.ShowLineNumbers = false
	editor := Editor{
		name: nameInput,
		desc: descInput,
		mode: EDIT_MODE_NAME,
	}
	return editor
}

func (e *Editor) ToggleTag(index int) {
	e.SetTag(index, !e.HasTag(index))
}

func (e Editor) HasTag(index int) bool {
	return (e.tags & (1 << index)) > 0
}

func (e *Editor) SetTag(index int, state bool) {
	if state {
		e.tags |= (1 << index)
	} else {
		e.tags &= ^(1 << index)
	}
}

func (e Editor) ViewTags() string {
	tags := []string{}
	for i, tag := range allTags {
		if !e.HasTag(i) {
			tag.color = lip.Color("#404040")
		}
		tags = append([]string{fmt.Sprintf(" %s", tag.View())}, tags...)
	}
	return lip.JoinHorizontal(lip.Top, tags...)
}

func (e Editor) LoadTask(task Task) Editor {
	e.task = task
	e.name.SetValue(task.name)
	e.desc.SetValue(task.description)
	e.tags = task.tags
	return e
}

func (e Editor) Init() tea.Cmd {
	return nil
}

func (e Editor) Update(msg tea.Msg, m model) (Editor, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.height = msg.Height
		e.desc.SetHeight(msg.Height-8)
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

		if e.mode == EDIT_MODE_NAME {
			nTags := len(allTags)
			for i := range nTags {
				key := fmt.Sprintf("f%d", nTags-i)
				if msg.String() != key {
					continue
				}
				e.ToggleTag(i)
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
	nameLabel := style.nameLabelStyle.Render(" Task ")
	name := e.name.View()
	desc := e.desc.View()
	tags := e.ViewTags()
	result := lip.JoinVertical(
		lip.Left,
		lip.JoinHorizontal(lip.Top, nameLabel, tags),
		"",
		name,
		"",
		style.descLabelStyle.Render(" Description "),
		"",
		desc,
	)
	result = editorStyle.containerStyle.MaxWidth(EDITOR_WIDTH-2).Render(result)
	// result = lip.JoinVertical(
	// 	lip.Right,
	// 	result,
	// 	projectTitleStyle.
	// 		MaxWidth(EDITOR_WIDTH).
	// 		MaxHeight(1).
	// 		Render(fmt.Sprintf(" %s ", m.title)),
	// )
	return result
}
