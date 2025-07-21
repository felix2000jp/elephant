package app

import (
	"elephant/internal/core"
	"elephant/internal/features/notes/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	State State

	listModel *list.Model
}

func NewModel() Model {
	repository := core.NewRepository("./elephant")
	listModel := list.NewModel(&repository)

	return Model{
		State:     ListState,
		listModel: &listModel,
	}
}

func (m *Model) Init() tea.Cmd {
	m.listModel.Init()
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	listCmd := m.listModel.Update(msg)

	return m, tea.Batch(listCmd)
}

func (m *Model) View() string {
	if m.State == ListState {
		return m.listModel.View()
	}
	return "Could not render application"
}
