package view

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
)

type customKeyMap struct {
	baseKeyMap   viewport.KeyMap
	editNote     key.Binding
	quitViewNote key.Binding
}

func newCustomKeyMap() customKeyMap {
	km := customKeyMap{
		editNote: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter/space", "edit note"),
		),
		quitViewNote: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back to list note"),
		),
	}

	km.baseKeyMap = viewport.DefaultKeyMap()
	return km
}
