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

	c.width = msg.Width - h
	c.height = msg.Height - v

	c.textarea.SetWidth(c.width)
	c.textarea.SetHeight(c.height)
	return nil
}
