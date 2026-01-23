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
}

func NewTask(mptr *model, name string) Task {
	return Task{
		name:  name,
		mptr:  mptr,
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
	for i, tag := range allTags {
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

func (t Task) View(m model, b Board, index int) string {
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
	// strName := fmt.Sprintf("%d %s", index+1, t.name)
	strName := t.name
	strDesc := t.description
	if len(strDesc) == 0 {
		strDesc = "..."
	}


	// render tags
	tags := ""
	taskTags := t.GetTags(m)
	for i := range len(taskTags) {
		tag := taskTags[(len(taskTags)-1)-i]
		// tag color muting in unfocussed mode
		if !b.selected {
			tag.color = lip.Color("#646464")
		}
		tags += fmt.Sprintf("%s ", tag.View())
	}

	// create spacing between name and tags (unused for now)
	tagsWidth := lip.Width(tags)
	descWidth := lip.Width(strDesc)
	nameWidth := lip.Width(strName)
	nameTrim := (b.width-4) - nameWidth
	descTagDistance := (b.width-4) - (tagsWidth+descWidth)
	if descTagDistance < 0 {
		trim := descTagDistance * -1
		trim += 3
		sliceBound := descWidth-trim
		if sliceBound >= 0 {
			strDesc = strDesc[:sliceBound]+"..."
		}
	}
	if nameTrim < 0 {
		trim := nameTrim * -1
		trim += 3
		sliceBound := nameWidth-trim
		if sliceBound >= 0 {
			strName = strName[:sliceBound]+"..."
		}
	}

	// render stylization
	desc := style.descriptionStyle.Render(strDesc)
	name := style.nameStyle.Render(strName)

	// assemble name and tags
	desc = lip.JoinHorizontal(
		lip.Top,
		tags,
		desc,
	)

	// put together full task
	// put together full task
	str_container := lip.JoinVertical(
		lip.Left,
		name,
		desc,
	)

	// render container
	result := style.containerStyle.MaxWidth(b.width-4).Render(str_container)

	return result
}
