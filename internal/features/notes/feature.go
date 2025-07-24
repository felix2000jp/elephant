package notes

import (
	"elephant/internal/core"
	"elephant/internal/features"
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

type Feature struct {
	State         State
	listComponent features.Component
	viewComponent features.Component
	editComponent features.Component
}

func NewFeature() Feature {
	repository := core.NewNoteRepository(".elephant/notes")
	listComponent := NewComponent(&repository)
	viewComponent := view.NewComponent(&repository)
	editComponent := edit.NewComponent(&repository)

	return Feature{
		State:         ListState,
		listComponent: &listComponent,
		viewComponent: &viewComponent,
		editComponent: &editComponent,
	}
}

func (m *Feature) Init() tea.Cmd {
	return tea.Batch(
		m.listComponent.Init(),
		m.viewComponent.Init(),
		m.editComponent.Init(),
	)
}

func (m *Feature) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if _, ok := msg.(list.NoteSelectedMsg); ok {
		m.State = ViewState
	}
	if _, ok := msg.(view.QuitNoteMarkdownMsg); ok {
		m.State = ListState
	}
	if _, ok := msg.(view.EditNoteContentMsg); ok {
		m.State = EditState
	}
	if _, ok := msg.(edit.QuitNoteTextareaMsg); ok {
		m.State = ViewState
	}

	switch m.State {
	case ListState:
		cmd = m.listComponent.ForegroundUpdate(msg)
		cmds = append(cmds, cmd)
	case ViewState:
		cmd = m.viewComponent.ForegroundUpdate(msg)
		cmds = append(cmds, cmd)
	case EditState:
		cmd = m.editComponent.ForegroundUpdate(msg)
		cmds = append(cmds, cmd)
	}

	cmd = m.listComponent.BackgroundUpdate(msg)
	cmds = append(cmds, cmd)

	cmd = m.viewComponent.BackgroundUpdate(msg)
	cmds = append(cmds, cmd)

	cmd = m.editComponent.BackgroundUpdate(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Feature) View() string {
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
