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

type UIStyle struct {
	titleStyle lip.Style
	helpStyle  lip.Style
}

var taskStyle          TaskStyle
var taskStyleSelected  TaskStyle
var taskStyleUnfocused TaskStyle

var boardStyle         BoardStyle
var boardStyleSelected BoardStyle

var uiStyle            UIStyle

func InitStyles() {
	taskStyle = TaskStyle{
		containerStyle: lip.NewStyle().
			PaddingBottom(1),
		nameStyle: lip.NewStyle().
			Bold(true).
			MaxHeight(1).
			Height(1).
			Foreground(lip.Color("#808080")),
		descriptionStyle: lip.NewStyle().
			Foreground(lip.Color("#404040")).
			MaxHeight(1).
			Height(1),
	}
	taskStyleSelected = TaskStyle{
		containerStyle: lip.NewStyle().
			PaddingBottom(1),
		nameStyle: lip.NewStyle().
			Bold(true).
			MaxHeight(1).
			Height(1).
			Background(lip.Color("7")).
			Foreground(lip.Color("0")),
		descriptionStyle: lip.NewStyle().
			Foreground(lip.Color("#808080")).
			MaxHeight(1).
			Height(1),
	}
	taskStyleUnfocused = TaskStyle{
		containerStyle: lip.NewStyle().
			PaddingBottom(1),
		nameStyle: lip.NewStyle().
			Bold(true).
			MaxHeight(1).
			Height(1).
			Foreground(lip.Color("#404040")),
		descriptionStyle: lip.NewStyle().
			Foreground(lip.Color("#2d2d2d")).
			MaxHeight(1).
			Height(1),
	}
	boardStyle = BoardStyle{
		containerStyle: lip.NewStyle().
			Padding(1, 2, 0).
			Border(lip.RoundedBorder()).
			BorderForeground(lip.Color("#404040")),
		titleStyle: lip.NewStyle().Foreground(lip.Color("#404040")),
	}
	boardStyleSelected = BoardStyle{
		containerStyle: lip.NewStyle().
			Padding(1, 2, 0).
			Border(lip.RoundedBorder()).
			BorderForeground(lip.Color("7")),
		titleStyle: lip.NewStyle().Bold(true),
	}
	uiStyle = UIStyle{
		helpStyle: lip.NewStyle().
			Foreground(lip.Color("#404040")),
		titleStyle: lip.NewStyle().
			Foreground(lip.Color("0")).
			Background(lip.Color("7")),
	}
}
