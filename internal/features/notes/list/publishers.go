package list

import (
	"elephant/internal/core"
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
)

type NotesLoadedMsg struct {
	Notes []core.Note
}

func (c *Component) PublishNotesLoaded() tea.Msg {
	notes, err := c.repository.GetAllNotes()
	if err != nil {
		slog.Error("failed to load notes", "error", err)
		return NotesLoadedMsg{}
	}

	return NotesLoadedMsg{Notes: notes}
}

type NoteSelectedMsg struct {
	Note core.Note
}

func (c *Component) PublishNotesLoadedMsg() tea.Msg {
	selectedItem := c.list.SelectedItem().(core.Note)
	return NoteSelectedMsg{Note: selectedItem}
}
