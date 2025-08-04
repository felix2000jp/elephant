package list

import (
	"github.com/charmbracelet/bubbles/key"
)

type componentKeyMap struct {
	addNote  key.Binding
	viewNote key.Binding
}

func newComponentKeyMap() componentKeyMap {
	km := componentKeyMap{
		addNote: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new note"),
		),
		viewNote: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter/space", "view note"),
		),
	}

	return km
}

func (a componentKeyMap) getListOfBindings() []key.Binding {
	return []key.Binding{
		a.addNote,
		a.viewNote,
	}
}
