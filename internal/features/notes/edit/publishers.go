package edit

import tea "github.com/charmbracelet/bubbletea"

type QuitNoteTextareaMsg struct{}

func (c *Component) PublishQuitNoteTextareaMsg() tea.Msg {
	return QuitNoteTextareaMsg{}
}
