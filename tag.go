package main

import (
	lip "github.com/charmbracelet/lipgloss"
)

const MAX_TAGS = 8 // nice number B^)

type Tag struct {
	icon  string // should only be a single char
	color lip.Color
}

var allTags []Tag

// tags added last have higher priority during sorting
func NewTag(icon string, color lip.Color) {
	if len(allTags) >= MAX_TAGS {
		return
	}
	tag := Tag{
		icon:  icon,
		color: color,
	}
	allTags = append([]Tag{tag}, allTags...)
}

func (t Tag) View() string {
	style := lip.NewStyle().Foreground(t.color)
	return style.Render(t.icon)
}
