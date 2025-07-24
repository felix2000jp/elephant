package features

import tea "github.com/charmbracelet/bubbletea"

type Component interface {
	Init() tea.Cmd
	BackgroundUpdate(msg tea.Msg) tea.Cmd
	ForegroundUpdate(msg tea.Msg) tea.Cmd
	View() string
}
