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
		containerStyle: lip.NewStyle().
			PaddingBottom(1),
		nameStyle: lip.NewStyle().
			Bold(true).
			Foreground(lip.Color("#808080")),
		descriptionStyle: lip.NewStyle().
			Foreground(lip.Color("#404040")),
	}
	taskStyleSelected = TaskStyle{
		containerStyle: lip.NewStyle().
			PaddingBottom(1),
		nameStyle: lip.NewStyle().
			Bold(true).
			Background(lip.Color("7")).
			Foreground(lip.Color("0")),
		descriptionStyle: lip.NewStyle().
			Foreground(lip.Color("#808080")),
	}
}
