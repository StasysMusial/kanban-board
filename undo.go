package main

import (
	"encoding/json"
)

type modelState struct {
	data         []byte
	boardCursor  int
	columnCursor int
}

const MAX_STACK_SIZE = 100

var undoStack []modelState
var redoStack []modelState

func AddUndoPoint(m model, resetRedo bool) {
	stackSize := len(undoStack)
	if stackSize >= MAX_STACK_SIZE {
		startIndex := stackSize - MAX_STACK_SIZE
		startIndex++
		undoStack = undoStack[startIndex:]
	}
	boardCursor := m.cursor
	columnCursor := m.columns[boardCursor].cursor
	state := modelState{
		data: ModelToJSON(m),
		boardCursor: boardCursor,
		columnCursor: columnCursor,
	}
	undoStack = append(undoStack, state)
	if resetRedo {
		redoStack = []modelState{}
	}
}

func AddRedoPoint(m model) {
	stackSize := len(redoStack)
	if stackSize >= MAX_STACK_SIZE {
		startIndex := stackSize - MAX_STACK_SIZE
		startIndex++
		redoStack = redoStack[startIndex:]
	}
	boardCursor := m.cursor
	columnCursor := m.columns[boardCursor].cursor
	state := modelState{
		data: ModelToJSON(m),
		boardCursor: boardCursor,
		columnCursor: columnCursor,
	}
	redoStack = append(redoStack, state)
}

func Undo(m *model, printMsg bool) {
	if len(undoStack) == 0 {
		return
	}
	AddRedoPoint(*m)
	state := undoStack[len(undoStack)-1]
	undoStack = undoStack[:len(undoStack)-1]
	var data modelSaveData
	json.Unmarshal(state.data, &data)
	m.LoadState(data)
	m.cursor = state.boardCursor
	m.columns[m.cursor].cursor = state.columnCursor
	if !m.columns[m.cursor].IsEmpty() {
		m.columns[m.cursor].tasks[state.columnCursor].selected = true
	}
	if printMsg {
		m.Print("Undo", msgColorInfo)
	}
}

func Redo(m *model) {
	if len(redoStack) == 0 {
		return
	}
	AddUndoPoint(*m, false)
	state := redoStack[len(redoStack)-1]
	redoStack = redoStack[:len(redoStack)-1]
	var data modelSaveData
	json.Unmarshal(state.data, &data)
	m.LoadState(data)
	m.cursor = state.boardCursor
	m.columns[m.cursor].cursor = state.columnCursor
	if !m.columns[m.cursor].IsEmpty() {
		m.columns[m.cursor].tasks[state.columnCursor].selected = true
	}
	m.Print("Redo", msgColorInfo)
}
