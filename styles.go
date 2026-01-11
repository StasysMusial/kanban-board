package main

import (
	lip "github.com/charmbracelet/lipgloss"
)

type TaskStyle struct {
	containerStyle   lip.Style
	nameStyle        lip.Style
	descriptionStyle lip.Style
	// tagsStyle        lip.Style
	// tags might need some kind of custom rendering
	// lets see...
}

var taskStyle         TaskStyle
var taskStyleSelected TaskStyle

func InitStyles() {
	taskStyle = TaskStyle{
		containerStyle: lip.NewStyle(),
		nameStyle: lip.NewStyle().
			Bold(true).
			Foreground(lip.Color("#808080")),
		descriptionStyle: lip.NewStyle().
			Foreground(lip.Color("#808080")),
	}
	taskStyleSelected = TaskStyle{
		containerStyle: lip.NewStyle().
			Background(lip.Color("#1c1c1c")),
		nameStyle: lip.NewStyle().
			Bold(true),
		descriptionStyle: lip.NewStyle().
			Foreground(lip.Color("#808080")),
	}
}
