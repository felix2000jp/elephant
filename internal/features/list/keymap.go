package list

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type customKeyMap struct {
	baseKeyMap list.KeyMap
	addNote    key.Binding
	viewNote   key.Binding
}

func newCustomKeyMap() customKeyMap {
	km := customKeyMap{
		addNote: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new note"),
		),
		viewNote: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter/space", "view note"),
		),
	}

	km.baseKeyMap = list.DefaultKeyMap()
	return km
}
