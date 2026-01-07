package main

import (
	"fmt"
	lip "github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

type List struct {
	width   int
	height  int
	cursor  int // what line is the cursor on
	scroll  int // how far are we scrolled down
	content []string
}

var listStyle lip.Style = lip.NewStyle().
	Border(lip.NormalBorder()).
	BorderForeground(lip.Color("#808080")).
	Padding(1, 3, 1, 2)

func NewList(w int, h int, c []string) List {
	list := List{
		width:   w,
		height:  h,
		cursor:  0,
		scroll:  0,
		content: c,
	}
	return list
}

func (l List) Update(msg tea.Msg) (List, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if l.cursor < len(l.content)-1 {
				l.cursor++
				if l.cursor > l.scroll + l.height {
					l.scroll++
				}
			}
		case "k", "up":
			if l.cursor > 0 {
				l.cursor--
				if l.cursor < l.scroll {
					l.scroll--
				}
			}
		}
	}
	return l, nil
}

func (l List) View() string {
	rendered := []string{}

	for i, line := range l.content {
		if i > l.scroll + l.height {
			continue
		}
		if i < l.scroll {
			continue
		}

		cursor := " "
		if l.cursor == i {
			cursor = ">"
		}

		rendered = append(rendered, fmt.Sprintf("%s %s", cursor, line))
	}

	s := lip.JoinVertical(lip.Left, rendered...)
	s = listStyle.Render(s)
	return s
}
