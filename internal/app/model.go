package app

import (
	"elephant/internal/core"
	"elephant/internal/features/notes/edit"
	"elephant/internal/features/notes/list"
	"elephant/internal/features/notes/view"
	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	ListState State = iota
	ViewState
	EditState
)

type Model struct {
	State         State
	listComponent *list.Component
	viewComponent *view.Component
	editComponent *edit.Component
}

func NewModel() Model {
	repository := core.NewNoteRepository(".elephant")
	listComponent := list.NewComponent(&repository)
	viewComponent := view.NewComponent(&repository)
	editComponent := edit.NewComponent(&repository)

	return Model{
		State:         ListState,
		listComponent: &listComponent,
		viewComponent: &viewComponent,
		editComponent: &editComponent,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.listComponent.Init(),
		m.viewComponent.Init(),
		m.editComponent.Init(),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Batch(
		m.listComponent.Update(msg),
		m.viewComponent.Update(msg),
		m.editComponent.Update(msg),
	)
}

func (m *Model) View() string {
	switch m.State {
	case ListState:
		return m.listComponent.View()
	case ViewState:
		return m.viewComponent.View()
	case EditState:
		return m.editComponent.View()
	default:
		return "Could not render application"
	}
}
