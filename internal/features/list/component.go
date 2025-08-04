package list

import (
	"elephant/internal/core"
	"elephant/internal/features/commands"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
)

type Component struct {
	width, height int
	list          list.Model
	keys          componentKeyMap
	repository    core.Repository
}

func NewComponent(repository core.Repository) Component {
	keys := newComponentKeyMap()
	itemList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	itemList.Title = "Elephant Notes"
	itemList.AdditionalFullHelpKeys = keys.getListOfBindings

	lc := Component{
		width:      itemList.Width(),
		height:     itemList.Height(),
		list:       itemList,
		keys:       keys,
		repository: repository,
	}

	return lc
}

func (lc *Component) Init() tea.Cmd {
	return func() tea.Msg {
		notes, err := lc.repository.GetAllNotes()
		if err != nil {
			slog.Error("failed to load notes", "error", err)
			return commands.ListNotesMsg{}
		}

		return commands.ListNotesMsg{Notes: notes}
	}
}

func (lc *Component) BackgroundUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := theme.Style.GetFrameSize()

		lc.width = msg.Width - h
		lc.height = msg.Height - v

		lc.list.SetSize(lc.width, lc.height)

	case commands.ListNotesMsg:
		notes := msg.Notes
		items := make([]list.Item, len(notes))

		for i, note := range notes {
			items[i] = note
		}

		lc.list.SetItems(items)

	case commands.QuitEditNoteMsg:
		items := lc.list.Items()
		updatedNote := msg.Note

		for i, item := range items {
			note := item.(core.Note)
			if note.FilePath() == updatedNote.FilePath() {
				items[i] = updatedNote
				break
			}
		}

		lc.list.SetItems(items)

	case commands.CreateNoteMsg:
		totalItems := append(lc.list.Items(), msg.Note)
		lc.list.SetItems(totalItems)
	}

	return nil
}

func (lc *Component) ForegroundUpdate(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && lc.list.FilterState() != list.Filtering {
		switch {
		case key.Matches(keyMsg, lc.keys.addNote):
			return func() tea.Msg {
				return commands.AddNoteMsg{}
			}
		case key.Matches(keyMsg, lc.keys.viewNote):
			selectedItem := lc.list.SelectedItem().(core.Note)
			return func() tea.Msg {
				return commands.ViewNoteMsg{Note: selectedItem}
			}
		}
	}

	var cmd tea.Cmd
	lc.list, cmd = lc.list.Update(msg)
	return cmd
}

func (lc *Component) View() string {
	listView := lc.list.View()
	return theme.Style.Width(lc.width).Height(lc.height).Render(listView)
}
