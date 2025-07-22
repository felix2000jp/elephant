package list

import (
	"elephant/internal/core"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Component struct {
	width, height int
	list          list.Model
	repository    core.Repository
}

func NewComponent(repository core.Repository) Component {
	itemList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	c := Component{
		list:       itemList,
		width:      itemList.Width(),
		height:     itemList.Height(),
		repository: repository,
	}

	return c
}

func (c *Component) Init() tea.Cmd {
	return c.HandleInit()
}

func (c *Component) BackgroundUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmd = c.HandleResizeWindow(msg)
	case NotesLoadedMsg:
		cmd = c.HandleNotesLoaded(msg)
	}

	return cmd
}

func (c *Component) ForegroundUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEnter && c.list.FilterState() != list.Filtering {
			return func() tea.Msg {
				selectedItem := c.list.SelectedItem().(core.Note)
				return NoteSelectedMsg{Note: selectedItem}
			}
		}
	}

	c.list, cmd = c.list.Update(msg)
	return cmd
}

func (c *Component) View() string {
	listView := c.list.View()
	return theme.Style.Width(c.width).Height(c.height).Render(listView)
}
