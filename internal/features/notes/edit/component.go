package edit

import (
	"elephant/internal/core"
	"elephant/internal/features/notes/list"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type Component struct {
	width, height int
	textarea      textarea.Model
	repository    core.Repository
}

func NewComponent(repository core.Repository) Component {
	ta := textarea.New()
	ta.Focus()
	ta.ShowLineNumbers = false

	c := Component{
		width:      ta.Width(),
		height:     ta.Height(),
		textarea:   ta,
		repository: repository,
	}

	return c
}

func (c *Component) Init() tea.Cmd {
	return c.HandleInit()
}

func (c *Component) BackgroundUpdate(msg tea.Msg) tea.Cmd {
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

	c.textarea, cmd = c.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *Component) ForegroundUpdate(msg tea.Msg) tea.Cmd {
	return nil
}

func (c *Component) View() string {
	listView := c.textarea.View()
	return theme.Style.Width(c.width).Height(c.height).Render(listView)
}
