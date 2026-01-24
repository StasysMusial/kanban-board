package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Version struct {
	major int
	minor int
	patch int
}

func (v Version) ToString() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

func (v *Version) FromString(str string) {
	parts := strings.Split(str, ".")
	for i, part := range parts {
		number, _ := strconv.Atoi(part)
		switch i {
		case 0:
			v.major = number
		case 1:
			v.minor = number
		case 2:
			v.patch = number
		}
	}
}

const CURRENT_VERSION = "1.0.0"
var appVersion Version

func InitVersion() {
	appVersion.FromString(CURRENT_VERSION)
}

func main() {
	InitVersion()
	InitIO()
	GenerateDefaultConfig()
	if !GenerateProjectConfig() {
		return
	}
	// set up main model
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
