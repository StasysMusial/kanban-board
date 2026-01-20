// this is purely visual and is only needed for the
// help section at the bottom of the screen
package main

// context data struct declaration
type KeyContextData struct {
	actions []string
	keys    []string
}

// context enum
type KeyContext int
const (
	KEY_CONTEXT_BOARDS KeyContext = iota
	KEY_CONTEXT_TASK
	KEY_CONTEXT_TEXTEDIT
	KEY_CONTEXTS
)

// action enums
type BoardAction int
const (
	BOARD_ACTION_SELECT BoardAction = iota
	BOARD_ACTION_MOVE_TASK
	BOARD_ACTION_GOTO
	BOARD_ACTION_EDIT_TASK
	BOARD_ACTION_ADD_TASK
	BOARD_ACTION_DELETE_TASK
	BOARD_ACTION_COPY_TASK
	BOARD_ACTION_PASTE_TASK
	BOARD_ACTION_SORT
	BOARD_ACTION_QUIT
	BOARD_ACTIONS
)

type TaskAction int
const (
	TASK_ACTION_SWITCH_FIELD TaskAction = iota
	TASK_ACTION_TOGGLE_TAG
	TASK_ACTION_INSERT
	TASK_ACTION_SUBMIT
	TASK_ACTION_CANCEL
	TASK_ACTIONS
)

type TextEditAction int
const (
	TEXTEDIT_ACTION_EXIT TextEditAction = iota
	TEXTEDIT_ACTIONS
)

// specific context data structs declaration
var keyContextBoards   KeyContextData
var keyContextTask     KeyContextData
var keyContextTextEdit KeyContextData

// context data arrays declaration
var boardActionNames []string
var boardActionKeys  []string

var taskActionNames []string
var taskActionKeys  []string

var textEditActionNames []string
var textEditActionKeys  []string

// main access point for key contexts
var keyContexts map[KeyContext]KeyContextData

// initialize key contexts
func InitKeyContexts() {
	keyContexts = make(map[KeyContext]KeyContextData)

	// create context data storage
	keyContextBoards = KeyContextData{
		actions: []string{},
		keys: []string{},
	}
	keyContextTask = KeyContextData{
		actions: []string{},
		keys: []string{},
	}
	keyContextTextEdit = KeyContextData{
		actions: []string{},
		keys: []string{},
	}

	// resize arrays
	for range BOARD_ACTIONS {
		keyContextBoards.actions = append(keyContextBoards.actions, "")
		keyContextBoards.keys = append(keyContextBoards.keys, "")
	}
	for range TASK_ACTIONS {
		keyContextTask.actions = append(keyContextTask.actions, "")
		keyContextTask.keys = append(keyContextTask.keys, "")
	}
	for range TEXTEDIT_ACTIONS {
		keyContextTextEdit.actions = append(keyContextTextEdit.actions, "")
		keyContextTextEdit.keys = append(keyContextTextEdit.keys, "")
	}

	// setup boards key and action labels
	keyContextBoards.keys[BOARD_ACTION_QUIT]            = "q"
	keyContextBoards.actions[BOARD_ACTION_QUIT]         = "quit"

	keyContextBoards.keys[BOARD_ACTION_ADD_TASK]        = "a"
	keyContextBoards.actions[BOARD_ACTION_ADD_TASK]     = "add"

	keyContextBoards.keys[BOARD_ACTION_SELECT]          = "h/j/k/l"
	keyContextBoards.actions[BOARD_ACTION_SELECT]       = "select"

	keyContextBoards.keys[BOARD_ACTION_MOVE_TASK]       = "H/J/K/L"
	keyContextBoards.actions[BOARD_ACTION_MOVE_TASK]    = "move"

	keyContextBoards.keys[BOARD_ACTION_GOTO]            = "g/G"
	keyContextBoards.actions[BOARD_ACTION_GOTO]         = "top/bottom"

	keyContextBoards.keys[BOARD_ACTION_EDIT_TASK]       = "enter"
	keyContextBoards.actions[BOARD_ACTION_EDIT_TASK]    = "edit"

	keyContextBoards.keys[BOARD_ACTION_DELETE_TASK]     = "x"
	keyContextBoards.actions[BOARD_ACTION_DELETE_TASK]  = "cut"

	keyContextBoards.keys[BOARD_ACTION_COPY_TASK]       = "y"
	keyContextBoards.actions[BOARD_ACTION_COPY_TASK]    = "yank"

	keyContextBoards.keys[BOARD_ACTION_PASTE_TASK]      = "p"
	keyContextBoards.actions[BOARD_ACTION_PASTE_TASK]   = "paste"

	keyContextBoards.keys[BOARD_ACTION_SORT]            = "s"
	keyContextBoards.actions[BOARD_ACTION_SORT]         = "sort"

	// setup task key and action labels
	keyContextTask.keys[TASK_ACTION_SWITCH_FIELD]    = "tab"
	keyContextTask.actions[TASK_ACTION_SWITCH_FIELD] = "switch field"

	keyContextTask.keys[TASK_ACTION_TOGGLE_TAG]      = "1/2/3/4"
	keyContextTask.actions[TASK_ACTION_TOGGLE_TAG]   = "switch field"

	keyContextTask.keys[TASK_ACTION_INSERT]          = "i"
	keyContextTask.actions[TASK_ACTION_INSERT]       = "insert"

	keyContextTask.keys[TASK_ACTION_CANCEL]          = "esc"
	keyContextTask.actions[TASK_ACTION_CANCEL]       = "cancel"

	keyContextTask.keys[TASK_ACTION_SUBMIT]          = "enter"
	keyContextTask.actions[TASK_ACTION_SUBMIT]       = "submit"

	// setup text edit key and action labels
	keyContextTextEdit.keys[TEXTEDIT_ACTION_EXIT]    = "esc"
	keyContextTextEdit.actions[TEXTEDIT_ACTION_EXIT] = "stop editing"

	// populate key context array
	keyContexts[KEY_CONTEXT_BOARDS] = keyContextBoards
	keyContexts[KEY_CONTEXT_TASK] = keyContextTask
	keyContexts[KEY_CONTEXT_TEXTEDIT] = keyContextTextEdit
}
