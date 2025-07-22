package view

import (
	"elephant/internal/core"
	"elephant/internal/features/notes/list"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"log/slog"
)

type Component struct {
	width, height int
	markdown      viewport.Model
	renderer      *glamour.TermRenderer
	repository    core.Repository
}

func NewComponent(repository core.Repository) Component {
	vp := viewport.New(0, 0)
	renderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(120))
	if err != nil {
		slog.Error("failed to initialize markdown renderer", "error", err)
		panic("failed to initialize markdown renderer")
	}

	c := Component{
		markdown:   vp,
		renderer:   renderer,
		width:      vp.Width,
		height:     vp.Height,
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
	case list.NoteSelectedMsg:
		cmd = c.HandleListNoteSelectedMsg(msg)
		cmds = append(cmds, cmd)
	}

	c.markdown, cmd = c.markdown.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *Component) View() string {
	markdownView := c.markdown.View()
	return theme.Style.Width(c.width).Height(c.height).Render(markdownView)
}
