package view

import (
	"github.com/charmbracelet/bubbles/key"
)

type componentKeyMap struct {
	editNote     key.Binding
	quitViewNote key.Binding
}

func newComponentKeyMap() componentKeyMap {
	km := componentKeyMap{
		editNote: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter/space", "edit note"),
		),
		quitViewNote: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back to list note"),
		),
	}

	return km
}
