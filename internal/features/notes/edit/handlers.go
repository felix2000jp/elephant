package edit

import (
	"elephant/internal/features/notes/list"
	"elephant/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
)

func (c *Component) HandleResizeWindow(msg tea.WindowSizeMsg) tea.Cmd {
	h, v := theme.Style.GetFrameSize()

	c.width = msg.Width - h
	c.height = msg.Height - v

	c.textarea.SetWidth(c.width)
	c.textarea.SetHeight(c.height)
	return nil
}

func (c *Component) HandleListNoteSelectedMsg(msg list.NoteSelectedMsg) tea.Cmd {
	c.textarea.SetValue(msg.Note.FileContent())
	return nil
}
