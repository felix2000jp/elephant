package list

import (
	"elephant/internal/core"
	"elephant/internal/features/commands"
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

func (m *mockRepository) CreateEmptyNote(filename string) (core.Note, error) {
	if m.err != nil {
		return core.Note{}, m.err
	}
	return core.NewNote(filename+".md", ""), nil
}

func TestNewListComponent(t *testing.T) {
	mockRepo := &mockRepository{}
	component := NewComponent(mockRepo)

	if component.repository != mockRepo {
		t.Error("Expected repository to be set correctly")
	}

	if component.list.Title != "Elephant Notes" {
		t.Errorf("Expected title to be 'Elephant Notes', got '%s'", component.list.Title)
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

		component := NewComponent(mockRepo)
		cmd := component.Init()

		if cmd == nil {
			t.Fatal("Expected Init to return a command")
		}

		msg := cmd()
		listMsg, ok := msg.(commands.ListNotesMsg)
		if !ok {
			t.Fatal("Expected ListNotesMsg from Init command")
		}

		if len(listMsg.Notes) != 2 {
			t.Errorf("Expected 2 notes, got %d", len(listMsg.Notes))
		}
	})

	t.Run("note loading fails because GetAllNotes errors", func(t *testing.T) {
		mockRepo := &mockRepository{
			err: errors.New("repository error"),
		}

		component := NewComponent(mockRepo)
		cmd := component.Init()

		if cmd == nil {
			t.Fatal("Expected Init to return a command")
		}

		msg := cmd()
		listMsg, ok := msg.(commands.ListNotesMsg)
		if !ok {
			t.Fatal("Expected ListNotesMsg from Init command")
		}

		if len(listMsg.Notes) != 0 {
			t.Error("Expected empty notes slice when repository fails")
		}
	})
}

func TestListComponentBackgroundUpdate(t *testing.T) {
	t.Run("ListNotesMsg sets items on the list", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		note1 := core.NewNote("note1.md", "# Note 1\nContent 1")
		note2 := core.NewNote("note2.md", "# Note 2\nContent 2")
		msg := commands.ListNotesMsg{Notes: []core.Note{note1, note2}}

		cmd := component.BackgroundUpdate(msg)

		if cmd != nil {
			t.Error("Expected backgroundUpdate to return nil for ListNotesMsg")
		}

		items := component.list.Items()
		if len(items) != 2 {
			t.Errorf("Expected 2 items in list, got %d", len(items))
		}
	})

	t.Run("QuitEditNoteMsg updates selected item on the list", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		note1 := core.NewNote("/path/note1.md", "# Note 1\nOriginal content")
		note2 := core.NewNote("/path/note2.md", "# Note 2\nContent 2")
		listMsg := commands.ListNotesMsg{Notes: []core.Note{note1, note2}}
		component.BackgroundUpdate(listMsg)

		updatedNote1 := core.NewNote("/path/note1.md", "# Note 1\nUpdated content")
		quitMsg := commands.QuitEditNoteMsg{Note: updatedNote1}

		cmd := component.BackgroundUpdate(quitMsg)

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

	t.Run("CreateNoteMsg adds note to existing list", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		note1 := core.NewNote("note1.md", "# Note 1\nContent 1")
		note2 := core.NewNote("note2.md", "# Note 2\nContent 2")
		listMsg := commands.ListNotesMsg{Notes: []core.Note{note1, note2}}
		component.BackgroundUpdate(listMsg)

		newNote := core.NewNote("newnote.md", "# New Note\nNew content")
		createMsg := commands.CreateNoteMsg{Note: newNote}

		cmd := component.BackgroundUpdate(createMsg)

		if cmd != nil {
			t.Error("Expected backgroundUpdate to return nil for CreateNoteMsg")
		}

		items := component.list.Items()
		if len(items) != 3 {
			t.Errorf("Expected 3 items in list, got %d", len(items))
		}

		addedNote := items[2].(core.Note)
		if addedNote.Title() != "newnote" {
			t.Errorf("Expected added note title 'newnote', got '%s'", addedNote.Title())
		}
	})
}

func TestListComponentForegroundUpdate(t *testing.T) {
	t.Run("Enter key creates ViewNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		note1 := core.NewNote("note1.md", "# Note 1\nContent 1")
		listMsg := commands.ListNotesMsg{Notes: []core.Note{note1}}
		component.BackgroundUpdate(listMsg)

		keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd := component.ForegroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected foregroundUpdate to return a command for Enter key")
		}

		msg := cmd()
		viewMsg, ok := msg.(commands.ViewNoteMsg)
		if !ok {
			t.Fatal("Expected ViewNoteMsg from Enter key command")
		}

		if viewMsg.Note.Title() != "note1" {
			t.Errorf("Expected note title 'note1', got '%s'", viewMsg.Note.Title())
		}
	})

	t.Run("Enter key during filtering does not create ViewNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		note1 := core.NewNote("note1.md", "# Note 1\nContent 1")
		listMsg := commands.ListNotesMsg{Notes: []core.Note{note1}}
		component.BackgroundUpdate(listMsg)

		filterKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
		component.ForegroundUpdate(filterKeyMsg)

		keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd := component.ForegroundUpdate(keyMsg)

		if cmd != nil {
			msg := cmd()
			if _, ok := msg.(commands.ViewNoteMsg); ok {
				t.Error("Expected no ViewNoteMsg during filtering mode")
			}
		}
	})

	t.Run("'n' key creates AddNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
		cmd := component.ForegroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected foregroundUpdate to return a command for 'n' key")
		}

		msg := cmd()
		_, ok := msg.(commands.AddNoteMsg)
		if !ok {
			t.Error("Expected AddNoteMsg from 'n' key command")
		}
	})

	t.Run("'n' key during filtering does not create AddNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		note1 := core.NewNote("note1.md", "# Note 1\nContent 1")
		listMsg := commands.ListNotesMsg{Notes: []core.Note{note1}}
		component.BackgroundUpdate(listMsg)

		filterKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
		component.ForegroundUpdate(filterKeyMsg)

		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
		cmd := component.ForegroundUpdate(keyMsg)

		if cmd != nil {
			msg := cmd()
			if _, ok := msg.(commands.AddNoteMsg); ok {
				t.Error("Expected no AddNoteMsg during filtering mode")
			}
		}
	})
}
