package app

import (
	"elephant/internal/core"
	"elephant/internal/features/notes/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var appStyle = lipgloss.NewStyle()

type Model struct {
	State State

	listModel *list.Model
}

func NewModel() Model {
	repository := core.NewNoteRepository(".elephant")
	listModel := list.NewModel(&repository)

	return Model{
		State:     ListState,
		listModel: &listModel,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.listModel.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	listCmd := m.listModel.Update(msg)

	return m, tea.Batch(listCmd)
}

func (m *Model) View() string {
	switch m.State {
	case ListState:
		return m.listModel.View(appStyle)
	default:
		return "Could not render application"
	}
}
