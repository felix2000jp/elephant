package edit

import (
	"elephant/internal/core"
	tea "github.com/charmbracelet/bubbletea"
)

type QuitNoteTextareaMsg struct {
	Note core.Note
}

func (c *Component) PublishQuitNoteTextareaMsg() tea.Msg {
	c.currentNote = core.NewNote(c.currentNote.FilePath(), c.textarea.Value())
	return QuitNoteTextareaMsg{Note: c.currentNote}
}
