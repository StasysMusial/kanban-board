package main

import (
	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type Board struct {
	width    int
	height   int

	title    string
	color    lip.Color

	selected bool
	scroll   int
	cursor   int
	tasks    []Task
}

func NewBoard(title string, color lip.Color) Board {
	return Board{
		title:    title,
		color:    color,

		selected: false,
		scroll:   0,
		cursor:   0,
		tasks:    []Task{},
	}
}

func (b Board) Init() tea.Cmd {
	return nil
}

func (b Board) Update(msg tea.Msg) (Board, tea.Cmd) {
	return b, nil
}

func (b Board) View() string {
	return ""
}
