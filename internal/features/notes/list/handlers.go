package list

import (
	"elephant/internal/ui/theme"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
)

func (m *Model) HandleInit() tea.Cmd {
	return func() tea.Msg {
		notes, err := m.repository.GetAllNotes()
		if err != nil {
			slog.Error("failed to load notes", "error", err)
			return nil
		}

		return NotesLoadedMsg{Notes: notes}
	}
}

func (m *Model) HandleNotesLoaded(msg NotesLoadedMsg) tea.Cmd {
	notes := msg.Notes
	items := make([]list.Item, len(notes))

	for i, note := range notes {
		items[i] = note
	}

	m.list.SetItems(items)
	return nil
}

func (m *Model) HandleResizeWindow(msg tea.WindowSizeMsg) tea.Cmd {
	h, v := theme.Style.GetFrameSize()

	m.Width = msg.Width - h
	m.Height = msg.Height - v

	m.list.SetSize(m.Width, m.Height)

	return nil
}
