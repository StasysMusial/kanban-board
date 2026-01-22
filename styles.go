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
	scrollerStyle  lip.Style
}

type HelpStyle struct {
	containerStyle lip.Style
	keyStyle       lip.Style
	actionStyle    lip.Style
}

type EditorStyle struct {
	containerStyle lip.Style
	scrollerStyle  lip.Style

	nameLabelStyle lip.Style
	nameFieldStyle lip.Style

	descLabelStyle lip.Style
	descFieldStyle lip.Style
}

var taskStyle           TaskStyle
var taskStyleSelected   TaskStyle
var taskStyleUnfocused  TaskStyle

var boardStyle          BoardStyle
var boardStyleSelected  BoardStyle

var helpStyle           HelpStyle

var editorStyle         EditorStyle
var editorStyleName     EditorStyle
var editorStyleDesc     EditorStyle

var msgInfoColor  lip.Color
var msgWarnColor  lip.Color
var msgErrorColor lip.Color

func InitStyles() {
	msgInfoColor = lip.Color("#678899")
	msgWarnColor = lip.Color("#999967")
	msgErrorColor = lip.Color("#996777")

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
			Foreground(lip.Color("#1c1c1c")),
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
			Padding(0, 2).
			Border(lip.RoundedBorder()).
			BorderForeground(lip.Color("#2d2d2d")),
		titleStyle: lip.NewStyle().Foreground(lip.Color("#404040")),
		scrollerStyle: lip.NewStyle().Foreground(lip.Color("#404040")),
	}
	boardStyleSelected = BoardStyle{
		containerStyle: lip.NewStyle().
			Padding(0, 2).
			Border(lip.RoundedBorder()).
			BorderForeground(lip.Color("#646464")),
		titleStyle: lip.NewStyle().Bold(true),
		scrollerStyle: lip.NewStyle().Foreground(lip.Color("#808080")),
	}
	helpStyle = HelpStyle{
		containerStyle: lip.NewStyle().
			MaxHeight(1).
			Height(1),
		keyStyle: lip.NewStyle().
			Bold(true).
			Foreground(lip.Color("#505050")),
		actionStyle: lip.NewStyle().
			Foreground(lip.Color("#404040")),
	}
	editorStyle = EditorStyle{
		containerStyle: lip.NewStyle().
			Padding(2, 4, 2),
		scrollerStyle: lip.NewStyle().Foreground(lip.Color("#404040")),
		nameLabelStyle: lip.NewStyle().
			Bold(true).
			Foreground(lip.Color("#1c1c1c")).
			Background(lip.Color("#505050")),
		nameFieldStyle: lip.NewStyle().
			Foreground(lip.Color("7")),
		descLabelStyle: lip.NewStyle().
			Bold(true).
			Foreground(lip.Color("#1c1c1c")).
			Background(lip.Color("#505050")),
		descFieldStyle: lip.NewStyle().
			Foreground(lip.Color("7")),
	}
	editorStyleName = EditorStyle{
		containerStyle: lip.NewStyle().
			Padding(2, 4, 2),
		scrollerStyle: lip.NewStyle().Foreground(lip.Color("#808080")),
		nameLabelStyle: lip.NewStyle().
			Bold(true).
			Foreground(lip.Color("#1c1c1c")).
			Background(lip.Color("7")),
		nameFieldStyle: lip.NewStyle().
			Foreground(lip.Color("7")),
		descLabelStyle: lip.NewStyle().
			Bold(true).
			Foreground(lip.Color("#1c1c1c")).
			Background(lip.Color("#505050")),
		descFieldStyle: lip.NewStyle().
			Foreground(lip.Color("7")),
	}
	editorStyleDesc = EditorStyle{
		containerStyle: lip.NewStyle().
			Padding(2, 4, 2),
		scrollerStyle: lip.NewStyle().Foreground(lip.Color("#808080")),
		nameLabelStyle: lip.NewStyle().
			Bold(true).
			Foreground(lip.Color("#1c1c1c")).
			Background(lip.Color("#505050")),
		nameFieldStyle: lip.NewStyle().
			Foreground(lip.Color("7")),
		descLabelStyle: lip.NewStyle().
			Bold(true).
			Foreground(lip.Color("#1c1c1c")).
			Background(lip.Color("7")),
		descFieldStyle: lip.NewStyle().
			Foreground(lip.Color("7")),
	}
}
