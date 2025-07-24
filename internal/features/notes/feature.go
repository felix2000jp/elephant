package notes

import (
	"elephant/internal/core"
	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	ListState State = iota
	ViewState
	EditState
)

type Feature struct {
	State         State
	listComponent *listComponent
	viewComponent *viewComponent
	editComponent *editComponent
}

func NewFeature() Feature {
	repository := core.NewNoteRepository(".elephant/notes")
	list := newListComponent(&repository)
	view := newViewComponent(&repository)
	edit := newEditComponent(&repository)

	return Feature{
		State:         ListState,
		listComponent: &list,
		viewComponent: &view,
		editComponent: &edit,
	}
}

func (m *Feature) Init() tea.Cmd {
	return tea.Batch(
		m.listComponent.init(),
		m.viewComponent.init(),
		m.editComponent.init(),
	)
}

func (m *Feature) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if _, ok := msg.(ViewNoteMsg); ok {
		m.State = ViewState
	}
	if _, ok := msg.(QuitViewNoteMsg); ok {
		m.State = ListState
	}
	if _, ok := msg.(EditNoteMsg); ok {
		m.State = EditState
	}
	if _, ok := msg.(QuitEditNoteMsg); ok {
		m.State = ViewState
	}

	switch m.State {
	case ListState:
		cmd = m.listComponent.foregroundUpdate(msg)
		cmds = append(cmds, cmd)
	case ViewState:
		cmd = m.viewComponent.foregroundUpdate(msg)
		cmds = append(cmds, cmd)
	case EditState:
		cmd = m.editComponent.foregroundUpdate(msg)
		cmds = append(cmds, cmd)
	}

	cmd = m.listComponent.backgroundUpdate(msg)
	cmds = append(cmds, cmd)

	cmd = m.viewComponent.backgroundUpdate(msg)
	cmds = append(cmds, cmd)

	cmd = m.editComponent.backgroundUpdate(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Feature) View() string {
	switch m.State {
	case ListState:
		return m.listComponent.view()
	case ViewState:
		return m.viewComponent.view()
	case EditState:
		return m.editComponent.view()
	default:
		return "Could not render application"
	}
}
