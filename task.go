package main

import (
	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type Task struct {
	width       int
	name        string
	description string
	tags        int
	selected    bool
	mptr        *model // not sure if this will work but i want to reference the main model while rendering the task here
}

func NewTask(mptr *model, name string) Task {
	return Task{
		mptr: mptr,
		name: name,
	}
}

func (t *Task) HasTag(tag_index int) bool {
	return t.tags & (1 << tag_index) == 1
}

func (t *Task) SetTag(tag_index int, state bool) {
	mask := 1 << tag_index
	if state {
		// set the bit to 1
		t.tags |= mask
	} else {
		// set the bit to 0
		// the "^" flips all the bits (apparently)
		t.tags &= ^mask
	}
}

func (t Task) Init() tea.Cmd {
	return nil
}

func (t Task) Update(msg tea.Msg) (Task, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.width = msg.Width
	}
	return t, cmd
}

func (t Task) View() string {
	var style TaskStyle
	if t.selected {
		style = taskStyleSelected
	} else {
		style = taskStyle
	}

	str_name := t.name
	str_desc := t.description

	// jesus christ what a stupid hack
	// (in case i forget) this is measuring the difference in char length
	// between the two lines of strings in the task (name, desc) and adding
	// space characters at the end of the shorter one so that it renders
	// the background properly. this can not be the intended way lol
	diff := len(str_name) - len(str_desc)
	if diff < 0 {
		for range diff * -1 {
			str_name += " "
		}
	} else if diff > 0 {
		for range diff {
			str_desc += " "
		}
	}

	name := style.nameStyle.Render(str_name)
	desc := style.descriptionStyle.Render(str_desc)

	container := lip.JoinVertical(
		lip.Left,
		name,
		desc,
	)

	return style.containerStyle.Render(container)
}
