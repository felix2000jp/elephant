package edit

import (
	"elephant/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
)

func (c *Component) HandleInit() tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}

func (c *Component) HandleResizeWindow(msg tea.WindowSizeMsg) tea.Cmd {
	h, v := theme.Style.GetFrameSize()

	c.Width = msg.Width - h
	c.Height = msg.Height - v

	c.textarea.SetWidth(c.Width)
	c.textarea.SetHeight(c.Height)
	return nil
}
