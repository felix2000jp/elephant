package notes

import (
	"elephant/internal/core"
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

func TestNewViewComponent(t *testing.T) {
	mockRepo := &mockRepository{}
	component := newViewComponent(mockRepo)

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
		component := newViewComponent(mockRepo)

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
	})

	t.Run("QuitEditNoteMsg updates current note and re-renders content", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := newViewComponent(mockRepo)

		originalNote := core.NewNote("test.md", "# Original Content\nOriginal text")
		viewMsg := ViewNoteMsg{Note: originalNote}
		component.backgroundUpdate(viewMsg)

		updatedNote := core.NewNote("test.md", "# Updated Content\nUpdated text")
		quitMsg := QuitEditNoteMsg{Note: updatedNote}

		cmd := component.backgroundUpdate(quitMsg)

		if cmd != nil {
			t.Error("Expected backgroundUpdate to return nil for QuitEditNoteMsg")
		}

		if component.currentNote.FileContent() != "# Updated Content\nUpdated text" {
			t.Error("Expected current note to be updated with new content")
		}
	})

	t.Run("handles markdown rendering errors gracefully", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := newViewComponent(mockRepo)

		note := core.NewNote("test.md", "```\nunclosed code block")
		msg := ViewNoteMsg{Note: note}

		cmd := component.backgroundUpdate(msg)

		if cmd != nil {
			t.Error("Expected backgroundUpdate to return nil even with render errors")
		}

		if component.currentNote.FileContent() != "```\nunclosed code block" {
			t.Error("Expected current note to be set even when rendering fails")
		}
	})
}

func TestViewComponentForegroundUpdate(t *testing.T) {
	t.Run("Escape key creates QuitViewNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := newViewComponent(mockRepo)

		keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
		cmd := component.foregroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected foregroundUpdate to return a command for Escape key")
		}

		msg := cmd()
		_, ok := msg.(QuitViewNoteMsg)
		if !ok {
			t.Error("Expected QuitViewNoteMsg from Escape key command")
		}
	})

	t.Run("Enter key creates EditNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := newViewComponent(mockRepo)

		keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd := component.foregroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected foregroundUpdate to return a command for Enter key")
		}

		msg := cmd()
		_, ok := msg.(EditNoteMsg)
		if !ok {
			t.Error("Expected EditNoteMsg from Enter key command")
		}
	})
}
