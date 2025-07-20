package components

import (
	"elephant/internal/core"
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

func (l *ListComponent) GetSelectedItem() core.Note {
	if item := l.list.SelectedItem(); item != nil {
		if note, ok := item.(core.Note); ok {
			return note
		}
	}

	return core.Note{}
}

func (l *ListComponent) SetItems(notes []core.Note) {
	var items []list.Item

	for _, note := range notes {
		items = append(items, note)
	}

	l.list.SetItems(items)
}
