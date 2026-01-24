package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"encoding/json"

	"github.com/BurntSushi/toml"
)

type (
	config struct {
		Title  string
		Color  string
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
	taskSaveData struct {
		Name string
		Desc string
		Tags int
	}
	boardSaveData struct {
		Tasks []taskSaveData
	}
	modelSaveData struct {
		Boards []boardSaveData
	}
)

const DEFAULT_CONF_DATA = `title = "DEFAULT"
color = "DEFAULT"
tags = [
	{ icon="󰫢", color="#ff4cc4" },
	{ icon="󰅩", color="#89d789" },
	{ icon="", color="#5c84d6" },
	{ icon="󰃣", color="#f5d33d" },
]
boards = [
	{
		name="IDEAS",
		icon="",
		color="#5d4cff",
	},
	{
		name="TO DO",
		icon="",
		color="#ff4cc4",
	},
	{
		name="IN PROGRESS",
		icon="",
		color="#ffcb4c",
	},
	{
		name="FINISHED",
		icon="",
		color="#bcff4c",
	},
]
`

const DEFAULT_CONF_DIR = "/.config/kanban/"
const DEFAULT_CONF     = "default_config.toml"

const PROJ_DIR  = "/.kanban/"
const PROJ_CONF = "config.toml"
const PROJ_DATA = "data.json"

var pathSeparator string
var userHome      string
var cwd           string

func InitIO() {
	pathSeparator = string(os.PathSeparator)
	userHome, _ = os.UserHomeDir()
	cwd, _ = os.Getwd()
}

func Path(input string) string {
	return strings.ReplaceAll(input, "/", pathSeparator)
}

func ModelToJSON(m model) []byte {
	var mData modelSaveData
	for _, board := range m.boards {
		var bData boardSaveData
		for _, task := range board.tasks {
			var tData taskSaveData
			tData.Name = task.name
			tData.Desc = task.description
			tData.Tags = task.tags
			bData.Tasks = append(bData.Tasks, tData)
		}
		mData.Boards = append(mData.Boards, bData)
	}
	json, _ := json.MarshalIndent(&mData, "", "	")
	return json
}

func GetConfig(path string) config {
	var c config
	toml.DecodeFile(path, &c)
	return c
}

func GenerateDefaultConfig() {
	path := GetDefaultConfigPath()
	// check if config already exists
	if _, err := os.Stat(path); err == nil {
		return
	}
	dirPath := GetDefaultConfigDir()
	os.MkdirAll(dirPath, os.ModeDir)

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	defer file.Close()
	file.WriteString(DEFAULT_CONF_DATA)
}

func GetDefaultConfigDir() string {
	return Path(userHome+DEFAULT_CONF_DIR)
}

func GetDefaultConfigPath() string {
	return Path(userHome+DEFAULT_CONF_DIR+DEFAULT_CONF)
}

func GetCwdConfigDir() string {
	return Path(cwd+PROJ_DIR)
}

func GetCwdConfigPath() string {
	return Path(cwd+PROJ_DIR+PROJ_CONF)
}

func ProjectPrompt() bool {
	fmt.Println("No kanban board detected. Initialize kanban board in this directory? (Y/n)")
	var response string
	fmt.Scanln(&response)
	switch response {
	case "Y", "y":
		return true
	case "N", "n":
		return false
	default:
		fmt.Println("Error: Unrecognized input")
		return ProjectPrompt()
	}
}

func GenerateProjectConfig() bool {
	path := GetCwdConfigPath()
	// check if config already exists
	if _, err := os.Stat(path); err == nil {
		return true
	} else if !ProjectPrompt() {
		return false
	}
	dirPath := GetCwdConfigDir()
	os.MkdirAll(dirPath, os.ModeDir)

	projConf, err := os.Create(path)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	defer projConf.Close()
	defaultConf, err := os.Open(GetDefaultConfigPath())
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	defer defaultConf.Close()
	io.Copy(projConf, defaultConf)

	return true
}

func GetCwdDataPath() string {
	return Path(cwd+PROJ_DIR+PROJ_DATA)
}

func WriteData(data []byte) {
	path := GetCwdDataPath()
	dirPath := GetCwdConfigDir()
	os.MkdirAll(dirPath, os.ModeDir)
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	defer file.Close()
	file.Write(data)
}

func ReadData() modelSaveData {
	var mData modelSaveData
	path := GetCwdDataPath()
	if _, err := os.Stat(path); err != nil {
		return mData
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	json.Unmarshal(data, &mData)
	return mData
}
