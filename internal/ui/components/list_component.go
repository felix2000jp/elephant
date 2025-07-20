package components

import (
	"elephant/internal/ui/theme"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListComponent struct {
	Width, Height int
	list          list.Model
}

func NewListComponent() *ListComponent {
	itemList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l := ListComponent{
		list:   itemList,
		Width:  itemList.Width(),
		Height: itemList.Height(),
	}

	return &l
}

func (l *ListComponent) Init() {
	l.list.Title = "Elephant Notes"
}

func (l *ListComponent) Update(msg tea.Msg) tea.Cmd {
	newNoteList, cmd := l.list.Update(msg)
	l.list = newNoteList

	return cmd
}

func (l *ListComponent) View() string {
	listView := l.list.View()
	return theme.Style.Width(l.Width).Height(l.Height).Render(listView)
}

func (l *ListComponent) ResizeWindow(msg tea.WindowSizeMsg) {
	h, v := theme.Style.GetFrameSize()

	l.Width = msg.Width - h
	l.Height = msg.Height - v

	l.list.SetSize(l.Width, l.Height)
}

func (l *ListComponent) GetSelectedItem() list.Item {
	if item := l.list.SelectedItem(); item != nil {
		return item
	}

	return nil
}

func (l *ListComponent) SetItems(items []list.Item) {
	l.list.SetItems(items)
}
