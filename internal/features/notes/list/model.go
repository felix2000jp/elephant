package list

import (
	"elephant/internal/core"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Width, Height int
	list          list.Model
	repository    core.Repository
}

func NewModel(repository core.Repository) Model {
	itemList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m := Model{
		list:       itemList,
		Width:      itemList.Width(),
		Height:     itemList.Height(),
		repository: repository,
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	return m.HandleInit()
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmd = m.HandleResizeWindow(msg)
		cmds = append(cmds, cmd)
	case NotesLoadedMsg:
		cmd = m.HandleNotesLoaded(msg)
		cmds = append(cmds, cmd)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Model) View(style lipgloss.Style) string {
	listView := m.list.View()
	return style.Width(m.Width).Height(m.Height).Render(listView)
}
