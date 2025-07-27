package notes

import (
	"elephant/internal/core"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"testing"
)

func TestNewAddComponent(t *testing.T) {
	mockRepo := &mockRepository{}
	component := newAddComponent(mockRepo)

	if component.repository != mockRepo {
		t.Error("Expected repository to be set correctly")
	}

	if !component.textInput.Focused() {
		t.Error("Expected text input to be focused on initialization")
	}

	if component.textInput.Placeholder != "Enter note filename (without .md)" {
		t.Errorf("Expected placeholder to be set correctly, got '%s'", component.textInput.Placeholder)
	}

	if component.width != 0 || component.height != 0 {
		t.Error("Expected initial dimensions to be 0")
	}
}

func TestAddComponentInit(t *testing.T) {
	mockRepo := &mockRepository{}
	component := newAddComponent(mockRepo)

	cmd := component.init()

	if cmd != nil {
		t.Error("Expected init to return nil")
	}
}

func TestAddComponentBackgroundUpdate(t *testing.T) {
	t.Run("CreateNoteMsg returns ViewNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := newAddComponent(mockRepo)

		note := core.NewNote("testnote.md", "")
		msg := CreateNoteMsg{Note: note}
		cmd := component.backgroundUpdate(msg)

		if cmd == nil {
			t.Fatal("Expected backgroundUpdate to return a command for CreateNoteMsg")
		}

		result := cmd()
		viewMsg, ok := result.(ViewNoteMsg)
		if !ok {
			t.Fatal("Expected ViewNoteMsg from CreateNoteMsg command")
		}

		if viewMsg.Note.Title() != "testnote" {
			t.Errorf("Expected note title 'testnote', got '%s'", viewMsg.Note.Title())
		}

		if viewMsg.Note.FilePath() != "testnote.md" {
			t.Errorf("Expected file path 'testnote.md', got '%s'", viewMsg.Note.FilePath())
		}
	})
}

func TestAddComponentForegroundUpdate(t *testing.T) {
	t.Run("Escape key creates QuitAddNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := newAddComponent(mockRepo)

		keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
		cmd := component.foregroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected foregroundUpdate to return a command for Escape key")
		}

		msg := cmd()
		_, ok := msg.(QuitAddNoteMsg)
		if !ok {
			t.Error("Expected QuitAddNoteMsg from Escape key command")
		}
	})

	t.Run("Enter key with filename creates note and CreateNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := newAddComponent(mockRepo)

		component.textInput.SetValue("mynote")

		keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd := component.foregroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected foregroundUpdate to return a command for Enter key with filename")
		}

		msg := cmd()
		createMsg, ok := msg.(CreateNoteMsg)
		if !ok {
			t.Fatal("Expected CreateNoteMsg from Enter key command")
		}

		if createMsg.Note.Title() != "mynote" {
			t.Errorf("Expected note title 'mynote', got '%s'", createMsg.Note.Title())
		}

		if createMsg.Note.FilePath() != "mynote.md" {
			t.Errorf("Expected note file path 'mynote.md', got '%s'", createMsg.Note.FilePath())
		}
	})

	t.Run("Enter key with filename handles repository error gracefully", func(t *testing.T) {
		mockRepo := &mockRepository{
			err: errors.New("repository error"),
		}
		component := newAddComponent(mockRepo)

		component.textInput.SetValue("mynote")

		keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd := component.foregroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected foregroundUpdate to return a command for Enter key with filename")
		}

		msg := cmd()
		if msg != nil {
			t.Error("Expected nil result when repository fails")
		}
	})

	t.Run("Enter key with empty filename returns nil", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := newAddComponent(mockRepo)

		component.textInput.SetValue("")

		keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd := component.foregroundUpdate(keyMsg)

		if cmd != nil {
			t.Error("Expected foregroundUpdate to return nil for Enter key with empty filename")
		}
	})
}
