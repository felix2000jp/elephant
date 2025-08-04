package add

import "github.com/charmbracelet/bubbles/key"

type componentKeyMap struct {
	createNote  key.Binding
	quitAddNote key.Binding
}

func newComponentKeyMap() componentKeyMap {
	km := componentKeyMap{
		createNote: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "create new note"),
		),
		quitAddNote: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back to list note"),
		),
	}

	return km
}
