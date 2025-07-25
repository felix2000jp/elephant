package notes

import (
	"elephant/internal/core"
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

func TestNewEditComponent(t *testing.T) {
	mockRepo := &mockRepository{}
	component := newEditComponent(mockRepo)

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
		component := newEditComponent(mockRepo)

		note := core.NewNote("test.md", "# Test Note\nThis is test content")
		msg := ViewNoteMsg{Note: note}

		cmd := component.backgroundUpdate(msg)

		if cmd != nil {
			t.Error("Expected backgroundUpdate to return nil for ViewNoteMsg")
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
		component := newEditComponent(mockRepo)

		originalNote := core.NewNote("test.md", "# Original Content")
		viewMsg := ViewNoteMsg{Note: originalNote}
		component.backgroundUpdate(viewMsg)

		component.textarea.SetValue("# Updated Content\nNew text added")

		keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
		cmd := component.foregroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected foregroundUpdate to return a command for Escape key")
		}

		msg := cmd()
		quitMsg, ok := msg.(QuitEditNoteMsg)
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
