package edit

import (
	"elephant/internal/features/notes/list"
	"elephant/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
)

func (c *Component) HandleInit() tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}

func (c *Component) HandleResizeWindow(msg tea.WindowSizeMsg) tea.Cmd {
	h, v := theme.Style.GetFrameSize()

	c.width = msg.Width - h
	c.height = msg.Height - v

	c.textarea.SetWidth(c.width)
	c.textarea.SetHeight(c.height)
	return nil
}

func (c *Component) HandleListNoteSelectedMsg(msg list.NoteSelectedMsg) tea.Cmd {
	note, err := c.repository.GetNoteByTitle(msg.NoteTitle)
	if err != nil {
		slog.Error("failed to load note", "error", err)
		c.textarea.SetValue("Could not render content.")
		return nil
	}

	c.textarea.SetValue(note.FileContent())
	return nil
}
