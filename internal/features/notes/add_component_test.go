package notes

import (
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
	t.Run("CreateNoteMsg creates note and returns ViewNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := newAddComponent(mockRepo)

		msg := CreateNoteMsg{Filename: "testnote"}
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

	t.Run("CreateNoteMsg handles repository error gracefully", func(t *testing.T) {
		mockRepo := &mockRepository{
			err: errors.New("repository error"),
		}
		component := newAddComponent(mockRepo)

		msg := CreateNoteMsg{Filename: "testnote"}
		cmd := component.backgroundUpdate(msg)

		if cmd == nil {
			t.Fatal("Expected backgroundUpdate to return a command for CreateNoteMsg")
		}

		result := cmd()
		if result != nil {
			t.Error("Expected nil result when repository fails")
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

	t.Run("Enter key with filename creates CreateNoteMsg", func(t *testing.T) {
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

		if createMsg.Filename != "mynote" {
			t.Errorf("Expected filename 'mynote', got '%s'", createMsg.Filename)
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
