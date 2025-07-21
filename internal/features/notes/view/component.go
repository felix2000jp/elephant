package view

import (
	"elephant/internal/core"
	"elephant/internal/features/notes"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type Component struct {
	Width, Height int
	markdown      viewport.Model
	renderer      *glamour.TermRenderer
	repository    core.Repository
}

func NewComponent(repository core.Repository) Component {
	vp := viewport.New(0, 0)
	c := Component{
		markdown:   vp,
		Width:      vp.Width,
		Height:     vp.Height,
		repository: repository,
	}

	return c
}

func (c *Component) Init() tea.Cmd {
	return c.HandleInit()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmd = c.HandleResizeWindow(msg)
		cmds = append(cmds, cmd)
	case notes.ViewNoteMsg:
		cmd = c.HandleViewNoteMsg(msg)
		cmds = append(cmds, cmd)
	}

	c.markdown, cmd = c.markdown.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *Component) View() string {
	markdownView := c.markdown.View()
	return theme.Style.Width(c.Width).Height(c.Height).Render(markdownView)
}
