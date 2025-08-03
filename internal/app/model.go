package app

import (
	"elephant/internal/features"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	notesFeature *features.NotesFeature
}

func NewModel() Model {
	notesFeature := features.NewFeature()

	return Model{
		notesFeature: &notesFeature,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.notesFeature.Init(),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Batch(
		m.notesFeature.Update(msg),
	)
}

func (m *Model) View() string {
	return m.notesFeature.View()
}
