package main

import ( "strconv"; "strings"; "fmt" )

// this number needs to be updated right before an update is released
const CURRENT_VERSION = "1.0.1"

var appVersion Version

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

func InitVersion() {
	appVersion.FromString(CURRENT_VERSION)
}
