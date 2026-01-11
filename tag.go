package main

import (
	"container/list"
	lip "github.com/charmbracelet/lipgloss"
)

type Tag struct {
	color lip.Color
	name  string
}

// the MAX_TAGS is there to prevent tag indicies that extend past
// the 64 bit length of the generic integer
const MAX_TAGS int = 16

func AddTag(tags *list.List, name string, color lip.Color) {
	if tags.Len() >= MAX_TAGS {
		// there should be either some kind of error message here
		// or a general way to prevent this
		return
	}
	tag := Tag{
		name: name,
		color: color,
	}
	tags.PushBack(tag)
}

func RemoveTag(tags *list.List, index int) {
	// not sure if these will actually be necessary,
	// guess we will see
	if index < 0 || index > tags.Len()-1 {
		return
	}
	if tags.Len() == 0 {
		return
	}

	// this code is still kind of confusing me.
	// it essentially allows to remove an element
	// of a linked list at a specified position
	// in that list
	element := tags.Front()
	for range index {
		element = element.Next()
	}
	tags.Remove(element)
}
