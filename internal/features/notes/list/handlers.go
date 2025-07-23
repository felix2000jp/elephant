package list

import (
	"elephant/internal/core"
	"elephant/internal/features/notes/edit"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (c *Component) HandleResizeWindow(msg tea.WindowSizeMsg) tea.Cmd {
	h, v := theme.Style.GetFrameSize()

	c.width = msg.Width - h
	c.height = msg.Height - v

	c.list.SetSize(c.width, c.height)
	return nil
}

func (c *Component) HandleNotesLoaded(msg NotesLoadedMsg) tea.Cmd {
	notes := msg.Notes
	items := make([]list.Item, len(notes))

	for i, note := range notes {
		items[i] = note
	}

	c.list.Title = "Elephant Notes"
	c.list.SetItems(items)
	return nil
}

func (c *Component) HandleQuitNoteTextareaMsg(msg edit.QuitNoteTextareaMsg) tea.Cmd {
	items := c.list.Items()
	updatedNote := msg.Note

	for i, item := range items {
		note := item.(core.Note)
		if note.FilePath() == updatedNote.FilePath() {
			items[i] = updatedNote
			break
		}
	}

	c.list.SetItems(items)
	return nil
}
