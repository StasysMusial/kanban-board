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

func (t Task) GetTags(m model) []Tag {
	tags := []Tag{}
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

func (t Task) View(m model, b Board) string {
	// select style
	var style TaskStyle
	if t.selected && b.selected {
		style = taskStyleSelected
	} else if !b.selected {
		style = taskStyleUnfocused
	} else {
		style = taskStyle
	}

	// get base strings
	str_name := t.name
	str_desc := t.description
	if len(str_desc) == 0 {
		str_desc = "..."
	}

	// render stylization
	name := style.nameStyle.Render(str_name)
	desc := style.descriptionStyle.Render(str_desc)

	// render tags
	tags := ""
	taskTags := t.GetTags(m)
	for i := range len(taskTags) {
		tag := taskTags[(len(taskTags)-1)-i]
		// tag color muting in unfocussed mode
		if !b.selected {
			tag.color = lip.Color("#646464")
		}
		tags += fmt.Sprintf(" %s", tag.View())
	}

	// create spacing between name and tags (unused for now)
	tags_width := lip.Width(tags)
	name_width := lip.Width(name)
	name_tag_distance := (b.width-4) - (tags_width+name_width)
	name_tag_spacing := ""
	for range name_tag_distance {
		name_tag_spacing += " "
	}

	// assemble name and tags
	title := lip.JoinHorizontal(
		lip.Top,
		name,
		// name_tag_spacing,
		tags,
	)

	// put together full task
	// put together full task
	str_container := lip.JoinVertical(
		lip.Left,
		title,
		desc,
	)

	// render container
	result := style.containerStyle.MaxWidth(b.width-4).Render(str_container)

	return result
}
