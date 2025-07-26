package notes

import (
	"elephant/internal/core"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
)

type listComponent struct {
	width, height int
	list          list.Model
	repository    core.Repository
}

func newListComponent(repository core.Repository) listComponent {
	itemList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	lc := listComponent{
		list:       itemList,
		width:      itemList.Width(),
		height:     itemList.Height(),
		repository: repository,
	}

	return lc
}

func (lc *listComponent) init() tea.Cmd {
	return func() tea.Msg {
		notes, err := lc.repository.GetAllNotes()
		if err != nil {
			slog.Error("failed to load notes", "error", err)
			return ListNotesMsg{}
		}

		return ListNotesMsg{Notes: notes}
	}
}

func (lc *listComponent) backgroundUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := theme.Style.GetFrameSize()

		lc.width = msg.Width - h
		lc.height = msg.Height - v

		lc.list.SetSize(lc.width, lc.height)

	case ListNotesMsg:
		notes := msg.Notes
		items := make([]list.Item, len(notes))

		for i, note := range notes {
			items[i] = note
		}

		lc.list.Title = "Elephant Notes"
		lc.list.SetItems(items)

	case QuitEditNoteMsg:
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
	}

	return nil
}

func (lc *listComponent) foregroundUpdate(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEnter && lc.list.FilterState() != list.Filtering {
			return func() tea.Msg {
				selectedItem := lc.list.SelectedItem().(core.Note)
				return ViewNoteMsg{Note: selectedItem}
			}
		}
		if keyMsg.Type == tea.KeyRunes && len(keyMsg.Runes) > 0 && keyMsg.Runes[0] == 'n' && lc.list.FilterState() != list.Filtering {
			return func() tea.Msg {
				return AddNoteMsg{}
			}
		}
	}

	var cmd tea.Cmd
	lc.list, cmd = lc.list.Update(msg)
	return cmd
}

func (lc *listComponent) view() string {
	listView := lc.list.View()
	return theme.Style.Width(lc.width).Height(lc.height).Render(listView)
}
