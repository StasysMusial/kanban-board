package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type Task struct {
	name        string
	description string
	tags        int
	selected    bool
	mptr        *model
	bptr        *Board
}

func NewTask(mptr *model, bptr *Board, name string) Task {
	return Task{
		name:  name,
		mptr:  mptr,
		bptr:  bptr,
	}
}

func (t *Task) SetTag(index int, state bool) {
	if state {
		t.tags |= (1 << index)
	} else {
		t.tags &= ^(1 << index)
	}
}

func (t Task) HasTag(index int) bool {
	return 0 < (t.tags & (1 << index))
}

func (t Task) GetTags() []Tag {
	tags := []Tag{}
	m := *t.mptr
	for i, tag := range m.tags {
		if t.HasTag(i) {
			tags = append(tags, tag)
		}
	}
	return tags
}

func (t Task) Init() tea.Cmd {
	return nil
}

func (t Task) Update(msg tea.Msg) (Task, tea.Cmd) {
	var cmd tea.Cmd
	return t, cmd
}

func (t Task) View() string {
	var style TaskStyle
	board := *t.bptr
	if t.selected && board.selected {
		style = taskStyleSelected
	} else if !board.selected {
		style = taskStyleUnfocused
	} else {
		style = taskStyle
	}

	str_name := t.name
	str_desc := t.description

	name := style.nameStyle.Render(str_name)
	desc := style.descriptionStyle.Render(str_desc)

	tags := ""
	for _, tag := range t.GetTags() {
		// tag color muting in unfocussed mode
		if !board.selected {
			tag.color = lip.Color("#646464")
		}
		tags += fmt.Sprintf(" %s", tag.View())
	}

	title := lip.JoinHorizontal(
		lip.Top,
		name,
		tags,
	)


	str_container := lip.JoinVertical(
		lip.Left,
		title,
		desc,
	)

	result := style.containerStyle.MaxWidth(board.width-4).Render(str_container)

	return result
}
