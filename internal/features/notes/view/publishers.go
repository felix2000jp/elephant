package view

import tea "github.com/charmbracelet/bubbletea"

type QuitNoteMarkdownMsg struct{}

func (c *Component) PublishQuitNoteMarkdownMsg() tea.Msg {
	return QuitNoteMarkdownMsg{}
}

type EditNoteContentMsg struct{}

func (c *Component) PublishEditNoteContentMsg() tea.Msg {
	return EditNoteContentMsg{}
}
