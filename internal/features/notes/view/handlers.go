package view

import (
	"elephant/internal/features/notes/list"
	"elephant/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
)

func (c *Component) HandleResizeWindow(msg tea.WindowSizeMsg) tea.Cmd {
	h, v := theme.Style.GetFrameSize()

	c.width = msg.Width - h
	c.height = msg.Height - v

	c.markdown.Width = c.width
	c.markdown.Height = c.height
	return nil
}

func (c *Component) HandleListNoteSelectedMsg(msg list.NoteSelectedMsg) tea.Cmd {
	content, err := c.renderer.Render(msg.Note.FileContent())
	if err != nil {
		slog.Error("failed to render markdown", "error", err)
		c.markdown.SetContent("Could not render content.")
		return nil
	}

	c.markdown.SetContent(content)
	return nil
}
