package view

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

func TestNewViewComponent(t *testing.T) {
	mockRepo := &mockRepository{}
	component := NewComponent(mockRepo)

	if component.repository != mockRepo {
		t.Error("Expected repository to be set correctly")
	}

	if component.renderer == nil {
		t.Error("Expected renderer to be initialized")
	}

	if component.width != 0 || component.height != 0 {
		t.Error("Expected initial dimensions to be 0")
	}
}

func TestViewComponentBackgroundUpdate(t *testing.T) {
	t.Run("ViewNoteMsg sets current note and renders content", func(t *testing.T) {
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
	})

	t.Run("QuitEditNoteMsg updates current note and re-renders content", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		originalNote := core.NewNote("test.md", "# Original Content\nOriginal text")
		viewMsg := commands.ViewNoteMsg{Note: originalNote}
		component.BackgroundUpdate(viewMsg)

		updatedNote := core.NewNote("test.md", "# Updated Content\nUpdated text")
		quitMsg := commands.QuitEditNoteMsg{Note: updatedNote}

		cmd := component.BackgroundUpdate(quitMsg)

		if cmd != nil {
			t.Error("Expected BackgroundUpdate to return nil for QuitEditNoteMsg")
		}

		if component.currentNote.FileContent() != "# Updated Content\nUpdated text" {
			t.Error("Expected current note to be updated with new content")
		}
	})

	t.Run("handles markdown rendering errors gracefully", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		note := core.NewNote("test.md", "```\nunclosed code block")
		msg := commands.ViewNoteMsg{Note: note}

		cmd := component.BackgroundUpdate(msg)

		if cmd != nil {
			t.Error("Expected BackgroundUpdate to return nil even with render errors")
		}

		if component.currentNote.FileContent() != "```\nunclosed code block" {
			t.Error("Expected current note to be set even when rendering fails")
		}
	})
}

func TestViewComponentForegroundUpdate(t *testing.T) {
	t.Run("Escape key creates QuitViewNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
		cmd := component.ForegroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected ForegroundUpdate to return a command for Escape key")
		}

		msg := cmd()
		_, ok := msg.(commands.QuitViewNoteMsg)
		if !ok {
			t.Error("Expected QuitViewNoteMsg from Escape key command")
		}
	})

	t.Run("Enter key creates EditNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd := component.ForegroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected ForegroundUpdate to return a command for Enter key")
		}

		msg := cmd()
		_, ok := msg.(commands.EditNoteMsg)
		if !ok {
			t.Error("Expected EditNoteMsg from Enter key command")
		}
	})
}
