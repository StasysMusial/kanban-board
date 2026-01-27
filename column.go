package main

import (
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type Column struct {
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

func NewColumn(title string, icon string, color lip.Color) Column {
	return Column{
		title:    title,
		icon:     icon,
		color:    color,

		selected: false,
		scroll:   0,
		cursor:   0,
		tasks:    []Task{},
	}
}

func (c *Column) ClampCursor() {
	c.cursor = max(c.cursor, 0)
	c.cursor = min(c.cursor, len(c.tasks)-1)
}

func (c *Column) IncCursor() {
	if c.cursor < len(c.tasks)-1 {
		c.cursor++
		if c.cursor >= c.shown_tasks + c.scroll {
			c.scroll++
		}
	}
}

func (c *Column) DecrCursor() {
	if c.cursor > 0 {
		c.cursor--
		if c.cursor < c.scroll {
			c.scroll--
		}
	}
}

func (c Column) IsEmpty() bool {
	return (len(c.tasks) == 0)
}

func (c *Column) AddTask(task Task, index int) {
	if index < 0 || index > len(c.tasks) {
		index = 0
	}
	tasks := []Task{}
	tasks = append(tasks, c.tasks[:index]...)
	tasks = append(tasks, task)
	tasks = append(tasks, c.tasks[index:]...)
	c.tasks = tasks
}

func (c *Column) RemoveTask(index int) {
	tasks := []Task{}
	tasks = append(c.tasks[:index], c.tasks[index+1:]...)
	c.tasks = tasks
	if len(c.tasks) == 0 {
		c.cursor = 0
	} else if c.cursor >= len(c.tasks) {
		c.cursor = len(c.tasks)-1
	}
	if c.cursor < c.scroll {
		c.scroll--
	}
}

// swaps task at index with index below
func (c *Column) SwapTask(index int) {
	index1 := index
	index2 := index-1

	if index1 >= len(c.tasks) {
		// index1 = 0
		return
	}
	if index2 < 0 {
		// index2 = len(b.tasks)-1
		return
	}

	t1 := c.tasks[index1]
	t2 := c.tasks[index2]
	c.tasks[index1] = t2
	c.tasks[index2] = t1
}

func (c *Column) MoveTaskUp() {
	c.SwapTask(c.cursor)
}

func (c *Column) MoveTaskDown() {
	c.SwapTask(c.cursor+1)
}

func (c *Column) SetSelected(selected bool) {
	c.selected = selected
	c.cursor = 0
	c.scroll = 0
}

func (c *Column) SetHeight(height int) {
	c.height = height
	c.shown_tasks = height / 3
}

func (c *Column) SortByTags(descending bool) {
	tasks := c.tasks
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
	c.tasks = tasks
}

func (c Column) Init() tea.Cmd {
	return nil
}

func (c Column) Update(msg tea.Msg, m model) (Column, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	b.height = msg.Height
	case tea.KeyMsg:
		if m.mode == MODE_NORMAL && c.selected && !c.IsEmpty() {
			switch msg.String() {
			case "j":
				c.IncCursor()
			case "k":
				c.DecrCursor()
			case "J":
				AddUndoPoint(m, true)
				c.MoveTaskDown()
				c.IncCursor()
			case "K":
				AddUndoPoint(m, true)
				c.MoveTaskUp()
				c.DecrCursor()
			case "g":
				c.cursor = 0
				c.scroll = 0
			case "G":
				c.cursor = len(c.tasks)-1
				c.scroll = c.cursor - c.shown_tasks
				c.scroll += 1
			case "s":
				AddUndoPoint(m, true)
				c.SortByTags(true)
			case "S":
				AddUndoPoint(m, true)
				c.SortByTags(false)
			}
		}
	}

	// Update tasks
	if c.selected {
		for i := range len(c.tasks) {
			var cmd tea.Cmd
			c.tasks[i], cmd = c.tasks[i].Update(msg)
			cmds = append(cmds, cmd)
			c.tasks[i].selected = (i == c.cursor)
		}
	}

	return c, tea.Batch(cmds...)
}

func (c Column) View(m model) string {
	var tasks []string
	for i, task := range c.tasks {
		if i < c.scroll {
			continue
		}
		tasks = append(tasks, task.View(m, c, i))
	}

	upScroller   := ""
	downScroller := ""

	if c.scroll > 0 {
		upScroller = "..."
	}

	if c.scroll < len(c.tasks) - c.shown_tasks {
		downScroller = "..."
	}

	var style ColumnStyle
	if c.selected && m.mode == MODE_NORMAL {
		style = columnStyleSelected
		// style.titleStyle = style.titleStyle
	} else {
		style = columnStyle
	}

	result := lip.JoinVertical(
		lip.Left,
		tasks...
	)

	result = lip.NewStyle().
		Width(c.width-4).
		Height(c.height).
		MaxHeight(c.height).
		Render(result)

	result = lip.JoinVertical(
		lip.Left,
		style.scrollerStyle.Width(c.width-4).AlignHorizontal(lip.Center).Render(upScroller),
		result,
		style.scrollerStyle.Width(c.width-4).AlignHorizontal(lip.Center).Render(downScroller),
	)

	icon := lip.NewStyle().Foreground(c.color).MaxHeight(1).Render(c.icon)
	counter := fmt.Sprintf("[%d]", len(c.tasks))
	counter = style.counterStyle.MaxHeight(1).Render(counter)
	title := style.titleStyle.MaxHeight(1).Render(c.title)

	header := lip.JoinHorizontal(lip.Top,
		icon,
		" ",
		title,
		" ",
		counter,
	)
	header = lip.NewStyle().
		Width(c.width).
		MaxWidth(c.width).
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
