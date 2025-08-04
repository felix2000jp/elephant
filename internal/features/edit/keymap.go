package edit

import "github.com/charmbracelet/bubbles/key"

type componentKeyMap struct {
	quitEditNote key.Binding
}

func newComponentKeyMap() componentKeyMap {
	km := componentKeyMap{
		quitEditNote: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back to view note"),
		),
	}

	return km
}
