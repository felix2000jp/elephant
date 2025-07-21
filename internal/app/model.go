package app

import (
	"elephant/internal/core"
	"elephant/internal/features/notes/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var appStyle = lipgloss.NewStyle()

type State int

const (
	ListState State = iota
	ViewState
	EditState
)

type Model struct {
	State         State
	listComponent *list.Component
}

func NewModel() Model {
	repository := core.NewNoteRepository(".elephant")
	listComponent := list.NewComponent(&repository)

	return Model{
		State:         ListState,
		listComponent: &listComponent,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.listComponent.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	listCmd := m.listComponent.Update(msg)

	return m, tea.Batch(listCmd)
}

func (m *Model) View() string {
	switch m.State {
	case ListState:
		return m.listComponent.View(appStyle)
	default:
		return "Could not render application"
	}
}
