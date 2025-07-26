package notes

import (
	"elephant/internal/core"
	tea "github.com/charmbracelet/bubbletea"
)

// ListNotesMsg - show the list of notes in the base path
type ListNotesMsg struct{ Notes []core.Note }

// ViewNoteMsg - select a single note to view and edit
type ViewNoteMsg struct{ Note core.Note }

// QuitViewNoteMsg - quit the view note state
type QuitViewNoteMsg struct{}

// EditNoteMsg - enter the edit state for the selected note
type EditNoteMsg struct{}

// QuitEditNoteMsg - quit the edit note state
type QuitEditNoteMsg struct{ Note core.Note }

// AddNoteMsg - enter the add note state
type AddNoteMsg struct{}

// CreateNoteMsg - create a new note with the given filename and enter edit state
type CreateNoteMsg struct{ Filename string }

// QuitAddNoteMsg - quit the add note state
type QuitAddNoteMsg struct{}

type State int

const (
	ListState State = iota
	ViewState
	EditState
	AddState
)

type Feature struct {
	State         State
	listComponent *listComponent
	viewComponent *viewComponent
	editComponent *editComponent
	addComponent  *addComponent
}

func NewFeature() Feature {
	repository := core.NewNoteRepository(".elephant/notes")
	list := newListComponent(&repository)
	view := newViewComponent(&repository)
	edit := newEditComponent(&repository)
	add := newAddComponent(&repository)

	return Feature{
		State:         ListState,
		listComponent: &list,
		viewComponent: &view,
		editComponent: &edit,
		addComponent:  &add,
	}
}

func (m *Feature) Init() tea.Cmd {
	return tea.Batch(
		m.listComponent.init(),
		m.viewComponent.init(),
		m.editComponent.init(),
		m.addComponent.init(),
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
	if _, ok := msg.(AddNoteMsg); ok {
		m.State = AddState
	}
	if _, ok := msg.(QuitAddNoteMsg); ok {
		m.State = ListState
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
	case AddState:
		cmd = m.addComponent.foregroundUpdate(msg)
		cmds = append(cmds, cmd)
	}

	cmd = m.listComponent.backgroundUpdate(msg)
	cmds = append(cmds, cmd)

	cmd = m.viewComponent.backgroundUpdate(msg)
	cmds = append(cmds, cmd)

	cmd = m.editComponent.backgroundUpdate(msg)
	cmds = append(cmds, cmd)

	cmd = m.addComponent.backgroundUpdate(msg)
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
	case AddState:
		return m.addComponent.view()
	default:
		return "Could not render application"
	}
}
