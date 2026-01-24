package main

import (
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type Board struct {
	width       int
	height      int

	title       string
	icon        string
	color       lip.Color

	selected    bool
	scroll      int
	cursor      int
	tasks       []Task

	shown_tasks int
}

func NewBoard(title string, icon string, color lip.Color) Board {
	return Board{
		title:    title,
		icon:     icon,
		color:    color,

		selected: false,
		scroll:   0,
		cursor:   0,
		tasks:    []Task{},
	}
}

func (b *Board) IncCursor() {
	if b.cursor < len(b.tasks)-1 {
		b.cursor++
		if b.cursor >= b.shown_tasks + b.scroll {
			b.scroll++
		}
	}
}

func (b *Board) DecrCursor() {
	if b.cursor > 0 {
		b.cursor--
		if b.cursor < b.scroll {
			b.scroll--
		}
	}
}

func (b Board) IsEmpty() bool {
	return (len(b.tasks) == 0)
}

func (b *Board) AddTask(task Task, index int) {
	if index < 0 || index > len(b.tasks) {
		index = 0
	}
	tasks := []Task{}
	tasks = append(tasks, b.tasks[:index]...)
	tasks = append(tasks, task)
	tasks = append(tasks, b.tasks[index:]...)
	b.tasks = tasks
}

func (b *Board) RemoveTask(index int) {
	tasks := []Task{}
	tasks = append(b.tasks[:index], b.tasks[index+1:]...)
	b.tasks = tasks
	if len(b.tasks) == 0 {
		b.cursor = 0
	} else if b.cursor >= len(b.tasks) {
		b.cursor = len(b.tasks)-1
	}
	if b.cursor < b.scroll {
		b.scroll--
	}
}

// swaps task at index with index below
func (b *Board) SwapTask(index int) {
	index1 := index
	index2 := index-1

	if index1 >= len(b.tasks) {
		// index1 = 0
		return
	}
	if index2 < 0 {
		// index2 = len(b.tasks)-1
		return
	}

	t1 := b.tasks[index1]
	t2 := b.tasks[index2]
	b.tasks[index1] = t2
	b.tasks[index2] = t1
}

func (b *Board) MoveTaskUp() {
	b.SwapTask(b.cursor)
}

func (b *Board) MoveTaskDown() {
	b.SwapTask(b.cursor+1)
}

func (b *Board) SetSelected(selected bool) {
	b.selected = selected
	b.cursor = 0
	b.scroll = 0
}

func (b *Board) SetHeight(height int) {
	b.height = height
	b.shown_tasks = height / 3
}

func (b *Board) SortByTags(descending bool) {
	tasks := b.tasks
	if descending {
		sort.Slice(tasks, func(i, j int) bool {
			t1 := tasks[i]
			t2 := tasks[j]
			return t1.tags > t2.tags
		})
	} else {
		sort.Slice(tasks, func(i, j int) bool {
			t1 := tasks[i]
			t2 := tasks[j]
			return t1.tags < t2.tags
		})
	}
	b.tasks = tasks
}

func (b Board) Init() tea.Cmd {
	return nil
}

func (b Board) Update(msg tea.Msg, m model) (Board, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	b.height = msg.Height
	case tea.KeyMsg:
		if m.mode == MODE_NORMAL && b.selected && !b.IsEmpty() {
			switch msg.String() {
			case "j":
				b.IncCursor()
			case "k":
				b.DecrCursor()
			case "J":
				b.MoveTaskDown()
				b.IncCursor()
			case "K":
				b.MoveTaskUp()
				b.DecrCursor()
			case "g":
				b.cursor = 0
				b.scroll = 0
			case "G":
				b.cursor = len(b.tasks)-1
				b.scroll = b.cursor - b.shown_tasks
				b.scroll += 1
			case "s":
				b.SortByTags(true)
			case "S":
				b.SortByTags(false)
			}
		}
	}

	// Update tasks
	if b.selected {
		for i := range len(b.tasks) {
			var cmd tea.Cmd
			b.tasks[i], cmd = b.tasks[i].Update(msg)
			cmds = append(cmds, cmd)
			b.tasks[i].selected = (i == b.cursor)
		}
	}

	return b, tea.Batch(cmds...)
}

func (b Board) View(m model) string {
	var tasks []string
	for i, task := range b.tasks {
		if i < b.scroll {
			continue
		}
		tasks = append(tasks, task.View(m, b, i))
	}

	upScroller   := ""
	downScroller := ""

	if b.scroll > 0 {
		upScroller = "..."
	}

	if b.scroll < len(b.tasks) - b.shown_tasks {
		downScroller = "..."
	}

	var style BoardStyle
	if b.selected && m.mode == MODE_NORMAL {
		style = boardStyleSelected
		// style.titleStyle = style.titleStyle
	} else {
		style = boardStyle
	}

	result := lip.JoinVertical(
		lip.Left,
		tasks...
	)

	result = lip.NewStyle().
		Width(b.width-4).
		Height(b.height).
		MaxHeight(b.height).
		Render(result)

	result = lip.JoinVertical(
		lip.Left,
		style.scrollerStyle.Width(b.width-4).AlignHorizontal(lip.Center).Render(upScroller),
		result,
		style.scrollerStyle.Width(b.width-4).AlignHorizontal(lip.Center).Render(downScroller),
	)

	icon := lip.NewStyle().Foreground(b.color).MaxHeight(1).Render(b.icon)
	counter := fmt.Sprintf("[%d]", len(b.tasks))
	counter = style.counterStyle.MaxHeight(1).Render(counter)
	title := style.titleStyle.MaxHeight(1).Render(b.title)

	header := lip.JoinHorizontal(lip.Top,
		icon,
		" ",
		title,
		" ",
		counter,
	)
	header = lip.NewStyle().
		Width(b.width).
		MaxWidth(b.width).
		Height(1).
		MaxHeight(1).
		AlignHorizontal(lip.Center).
		Render(header)

	result = style.containerStyle.Render(result)

	result = lip.JoinVertical(
		lip.Left,
		header,
		result,
	)

	return result
}
