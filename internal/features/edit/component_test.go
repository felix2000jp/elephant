package edit

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

func TestNewEditComponent(t *testing.T) {
	mockRepo := &mockRepository{}
	component := NewComponent(mockRepo)

	if component.repository != mockRepo {
		t.Error("Expected repository to be set correctly")
	}

	if component.textarea.Focused() != true {
		t.Error("Expected textarea to be focused on initialization")
	}

	if component.width == 0 || component.height == 0 {
		t.Error("Expected initial dimensions to be set from textarea")
	}
}

func TestEditComponentBackgroundUpdate(t *testing.T) {
	t.Run("ViewNoteMsg sets current note and textarea content", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		note := core.NewNote("test.md", "# Test Note\nThis is test content")
		msg := commands.ViewNoteMsg{Note: note}

		cmd := component.BackgroundUpdate(msg)

		if cmd != nil {
			t.Error("Expected BackgroundUpdate to return nil for ViewNoteMsg")
		}

		if component.currentNote.Title() != "test" {
			t.Errorf("Expected current note title to be 'test', got '%s'", component.currentNote.Title())
		}

		if component.currentNote.FileContent() != "# Test Note\nThis is test content" {
			t.Error("Expected current note content to be set correctly")
		}

		if component.textarea.Value() != "# Test Note\nThis is test content" {
			t.Error("Expected textarea to be populated with note content")
		}
	})
}

func TestEditComponentForegroundUpdate(t *testing.T) {
	t.Run("Escape key creates QuitEditNoteMsg with updated content", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		originalNote := core.NewNote("test.md", "# Original Content")
		viewMsg := commands.ViewNoteMsg{Note: originalNote}
		component.BackgroundUpdate(viewMsg)

		component.textarea.SetValue("# Updated Content\nNew text added")

		keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
		cmd := component.ForegroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected ForegroundUpdate to return a command for Escape key")
		}

		msg := cmd()
		quitMsg, ok := msg.(commands.QuitEditNoteMsg)
		if !ok {
			t.Fatal("Expected QuitEditNoteMsg from Escape key command")
		}

		if quitMsg.Note.FileContent() != "# Updated Content\nNew text added" {
			t.Error("Expected QuitEditNoteMsg to contain updated content from textarea")
		}

		if quitMsg.Note.FilePath() != "test.md" {
			t.Error("Expected QuitEditNoteMsg to preserve original file path")
		}
	})
}
