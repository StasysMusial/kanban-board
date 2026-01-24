package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type (
	config struct {
		Tags   []tagData
		Boards []boardData
	}
	tagData struct {
		Icon  string
		Color string
	}
	boardData struct {
		Name  string
		Icon  string
		Color string
	}
)

const DEFAULT_CONF_DATA = `tags = [
	{ icon="󰫢", color="#ff4cc4" },
	{ icon="󰅩", color="#89d789" },
	{ icon="", color="#5c84d6" },
	{ icon="󰃣", color="#f5d33d" },
]
boards = [
	{
		name="IDEAS",
		icon="",
		color="7",
	},
	{
		name="TO DO",
		icon="",
		color="7",
	},
	{
		name="IN PROGRESS",
		icon="",
		color="7",
	},
	{
		name="FINISHED",
		icon="",
		color="7",
	},
]
`

const DEFAULT_CONF_DIR = "/.config/kanban/"
const DEFAULT_CONF     = "default_config.toml"

const PROJ_DIR  = "/.kanban/"
const PROJ_CONF = PROJ_DIR + "config.toml"
const PROJ_DATA = PROJ_DIR + "data.json"

var pathSeparator string
var userHome      string

func InitIO() {
	pathSeparator = string(os.PathSeparator)
	userHome, _ = os.UserHomeDir()
}

func Path(input string) string {
	return strings.ReplaceAll(input, "/", pathSeparator)
}

func GetConfig(path string) config {
	var c config
	toml.DecodeFile(path, &c)
	return c
}

func GenerateDefaultConfig() {
	dirPath := GetDefaultConfigDir()
	os.MkdirAll(dirPath, os.ModeDir)
	path := GetDefaultConfigPath()
	// check if config already exists
	if _, err := os.Stat(path); err == nil {
		return
	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	file.WriteString(DEFAULT_CONF_DATA)
	defer file.Close()
}

func GetDefaultConfigDir() string {
	return Path(userHome+DEFAULT_CONF_DIR)
}

func GetDefaultConfigPath() string {
	return Path(userHome+DEFAULT_CONF_DIR+DEFAULT_CONF)
}
