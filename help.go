package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type Help struct {
	context_data KeyContextData
}

func (h *Help) SetKeyContext(context KeyContext) {
	h.context_data = keyContexts[context]
}

func (h Help) Init() tea.Cmd {
	return nil
}

func (h Help) Update(msg tea.Msg) (Help, tea.Cmd) {
	return h, nil
}

func (h Help) View(m model, maxWidth int) string {
	mappings := []string{}
	for i := range len(h.context_data.actions) {
		action := h.context_data.actions[i]
		key := h.context_data.keys[i]
		// action = helpStyle.actionStyle.Render(fmt.Sprintf("∙ %s", action))
		action = helpStyle.actionStyle.Render(action)
		key = helpStyle.keyStyle.Render(key)
		mapping := fmt.Sprintf("  %s %s", key, action)
		mappings = append(mappings, mapping)
	}

	mode := ""
	var modeColor lip.Color
	modelMode := m.mode
	editMode := m.editor.mode

	switch modelMode {
	case MODE_NORMAL:
		mode = "BOARD"
		modeColor = modeColorBoard
	case MODE_EDIT:
		switch editMode {
		case EDIT_MODE_NAME:
			mode = "TASK"
			modeColor = modeColorTask
		case EDIT_MODE_DESC:
			mode = "DESC"
			modeColor = modeColorDesc
		}
	}

	modeStyle := lip.NewStyle().
		Bold(true).
		Foreground(lip.Color("#1c1c1c")).
		Background(modeColor)
	mode = modeStyle.Render(fmt.Sprintf(" %s ", mode))

	result := lip.JoinHorizontal(lip.Top, mappings...)
	result = lip.JoinHorizontal(lip.Top, mode, result)
	result = helpStyle.containerStyle.
		MaxWidth(maxWidth).
		// AlignHorizontal(lip.Center).
		Render(result)

	return result
}
