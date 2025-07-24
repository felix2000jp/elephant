package features

import tea "github.com/charmbracelet/bubbletea"

type Feature interface {
	Init() tea.Cmd
	Update(msg tea.Msg) tea.Cmd
	View() string
}
