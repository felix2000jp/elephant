package notes

import (
	"elephant/internal/core"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

type mockRepository struct {
	notes []core.Note
	err   error
}

func (m *mockRepository) GetAllNotes() ([]core.Note, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.notes, nil
}

func (m *mockRepository) GetNoteByTitle(title string) (core.Note, error) {
	if m.err != nil {
		return core.Note{}, m.err
	}
	for _, note := range m.notes {
		if note.Title() == title {
			return note, nil
		}
	}
	return core.Note{}, errors.New("note not found")
}

func (m *mockRepository) SaveNote(_ core.Note) error {
	return m.err
}

func TestNewListComponent(t *testing.T) {
	mockRepo := &mockRepository{}
	component := newListComponent(mockRepo)

	if component.repository != mockRepo {
		t.Error("Expected repository to be set correctly")
	}

	if component.width != 0 || component.height != 0 {
		t.Error("Expected initial dimensions to be 0")
	}
}

func TestListComponentInit(t *testing.T) {
	t.Run("note loading is successful", func(t *testing.T) {
		note1 := core.NewNote("note1.md", "# Note 1\nContent 1")
		note2 := core.NewNote("note2.md", "# Note 2\nContent 2")
		mockRepo := &mockRepository{
			notes: []core.Note{note1, note2},
		}

		component := newListComponent(mockRepo)
		cmd := component.init()

		if cmd == nil {
			t.Fatal("Expected init to return a command")
		}

		msg := cmd()
		listMsg, ok := msg.(ListNotesMsg)
		if !ok {
			t.Fatal("Expected ListNotesMsg from init command")
		}

		if len(listMsg.Notes) != 2 {
			t.Errorf("Expected 2 notes, got %d", len(listMsg.Notes))
		}
	})

	t.Run("note loading fails because GetAllNotes errors", func(t *testing.T) {
		mockRepo := &mockRepository{
			err: errors.New("repository error"),
		}

		component := newListComponent(mockRepo)
		cmd := component.init()

		if cmd == nil {
			t.Fatal("Expected init to return a command")
		}

		msg := cmd()
		listMsg, ok := msg.(ListNotesMsg)
		if !ok {
			t.Fatal("Expected ListNotesMsg from init command")
		}

		if len(listMsg.Notes) != 0 {
			t.Error("Expected empty notes slice when repository fails")
		}
	})
}

func TestListComponentBackgroundUpdate(t *testing.T) {
	t.Run("ListNotesMsg sets items on the list", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := newListComponent(mockRepo)

		note1 := core.NewNote("note1.md", "# Note 1\nContent 1")
		note2 := core.NewNote("note2.md", "# Note 2\nContent 2")
		msg := ListNotesMsg{Notes: []core.Note{note1, note2}}

		cmd := component.backgroundUpdate(msg)

		if cmd != nil {
			t.Error("Expected backgroundUpdate to return nil for ListNotesMsg")
		}

		if component.list.Title != "Elephant Notes" {
			t.Errorf("Expected title to be 'Elephant Notes', got '%s'", component.list.Title)
		}

		items := component.list.Items()
		if len(items) != 2 {
			t.Errorf("Expected 2 items in list, got %d", len(items))
		}
	})

	t.Run("QuitEditNoteMsg updates selected item on the list", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := newListComponent(mockRepo)

		note1 := core.NewNote("/path/note1.md", "# Note 1\nOriginal content")
		note2 := core.NewNote("/path/note2.md", "# Note 2\nContent 2")
		listMsg := ListNotesMsg{Notes: []core.Note{note1, note2}}
		component.backgroundUpdate(listMsg)

		updatedNote1 := core.NewNote("/path/note1.md", "# Note 1\nUpdated content")
		quitMsg := QuitEditNoteMsg{Note: updatedNote1}

		cmd := component.backgroundUpdate(quitMsg)

		if cmd != nil {
			t.Error("Expected backgroundUpdate to return nil for QuitEditNoteMsg")
		}

		items := component.list.Items()
		if len(items) != 2 {
			t.Errorf("Expected 2 items in list, got %d", len(items))
		}

		updatedItem := items[0].(core.Note)
		if updatedItem.FileContent() != "# Note 1\nUpdated content" {
			t.Error("Expected first note to be updated with new content")
		}

		secondItem := items[1].(core.Note)
		if secondItem.FileContent() != "# Note 2\nContent 2" {
			t.Error("Expected second note to remain unchanged")
		}
	})
}

func TestListComponentForegroundUpdate(t *testing.T) {
	t.Run("Enter key creates ViewNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := newListComponent(mockRepo)

		note1 := core.NewNote("note1.md", "# Note 1\nContent 1")
		listMsg := ListNotesMsg{Notes: []core.Note{note1}}
		component.backgroundUpdate(listMsg)

		keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd := component.foregroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected foregroundUpdate to return a command for Enter key")
		}

		msg := cmd()
		viewMsg, ok := msg.(ViewNoteMsg)
		if !ok {
			t.Fatal("Expected ViewNoteMsg from Enter key command")
		}

		if viewMsg.Note.Title() != "note1" {
			t.Errorf("Expected note title 'note1', got '%s'", viewMsg.Note.Title())
		}
	})

	t.Run("Enter key during filtering does not create ViewNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := newListComponent(mockRepo)

		note1 := core.NewNote("note1.md", "# Note 1\nContent 1")
		listMsg := ListNotesMsg{Notes: []core.Note{note1}}
		component.backgroundUpdate(listMsg)

		filterKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
		component.foregroundUpdate(filterKeyMsg)

		keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd := component.foregroundUpdate(keyMsg)

		if cmd != nil {
			msg := cmd()
			if _, ok := msg.(ViewNoteMsg); ok {
				t.Error("Expected no ViewNoteMsg during filtering mode")
			}
		}
	})
}
