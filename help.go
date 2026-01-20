package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type Help struct {
	width        int
	context_data KeyContextData
}

func (h *Help) SetKeyContext(context KeyContext) {
	h.context_data = keyContexts[context]
}

func (h Help) Init() tea.Cmd {
	return nil
}

func (h Help) Update(msg tea.Msg) (Help, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.Width
	}
	return h, nil
}

func (h Help) View() string {
	mappings := []string{}
	for i := range len(h.context_data.actions) {
		action := h.context_data.actions[i]
		key := h.context_data.keys[i]
		// action = helpStyle.actionStyle.Render(fmt.Sprintf("∙ %s", action))
		action = helpStyle.actionStyle.Render(action)
		key = helpStyle.keyStyle.Render(key)
		mapping := fmt.Sprintf(" %s %s  ", key, action)
		mappings = append(mappings, mapping)
	}

	result := lip.JoinHorizontal(lip.Top, mappings...)
	result = helpStyle.containerStyle.
		Width(h.width).
		AlignHorizontal(lip.Center).
		Render(result)

	return result
}
