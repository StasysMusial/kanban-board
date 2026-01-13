package main

import (
	lip "github.com/charmbracelet/lipgloss"
)

type TaskStyle struct {
	containerStyle   lip.Style
	nameStyle        lip.Style
	descriptionStyle lip.Style
}

type BoardStyle struct {
	containerStyle lip.Style
	titleStyle     lip.Style
}

var taskStyle          TaskStyle
var taskStyleSelected  TaskStyle

var boardStyle         BoardStyle
var boardStyleSelected BoardStyle

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
	boardStyle = BoardStyle{
		containerStyle: lip.NewStyle().Padding(2, 3, 1),
		titleStyle: lip.NewStyle().Bold(true),
	}
	boardStyleSelected = BoardStyle{
		containerStyle: lip.NewStyle().Padding(1, 2, 0).Border(lip.NormalBorder()),
		titleStyle: lip.NewStyle().Bold(true),
	}
}
