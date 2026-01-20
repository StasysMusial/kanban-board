package main

import (
	lip "github.com/charmbracelet/lipgloss"
)

type Tag struct {
	icon  string // should only be a single char
	color lip.Color
}

// tags added last have higher priority during sorting
func NewTag(icon string, color lip.Color) Tag {
	return Tag{
		icon:  icon,
		color: color,
	}
}

func (t Tag) View() string {
	style := lip.NewStyle().Foreground(t.color)
	return style.Render(t.icon)
}
