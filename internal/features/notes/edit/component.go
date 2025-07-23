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
	return nil
}

func (c *Component) BackgroundUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmd = c.HandleResizeWindow(msg)
	case list.NoteSelectedMsg:
		cmd = c.HandleListNoteSelectedMsg(msg)
	}

	return cmd
}

func (c *Component) ForegroundUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEsc {
			return c.PublishQuitNoteTextareaMsg
		}
	}

	c.textarea, cmd = c.textarea.Update(msg)
	return cmd
}

func (c *Component) View() string {
	listView := c.textarea.View()
	return theme.Style.Width(c.width).Height(c.height).Render(listView)
}
