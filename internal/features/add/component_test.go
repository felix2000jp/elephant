package add

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

func TestNewAddComponent(t *testing.T) {
	mockRepo := &mockRepository{}
	component := NewComponent(mockRepo)

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
	component := NewComponent(mockRepo)

	cmd := component.Init()

	if cmd != nil {
		t.Error("Expected Init to return nil")
	}
}

func TestAddComponentBackgroundUpdate(t *testing.T) {
	t.Run("CreateNoteMsg returns ViewNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		note := core.NewNote("testnote.md", "")
		msg := commands.CreateNoteMsg{Note: note}
		cmd := component.BackgroundUpdate(msg)

		if cmd == nil {
			t.Fatal("Expected BackgroundUpdate to return a command for CreateNoteMsg")
		}

		result := cmd()
		viewMsg, ok := result.(commands.ViewNoteMsg)
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
		component := NewComponent(mockRepo)

		keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
		cmd := component.ForegroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected ForegroundUpdate to return a command for Escape key")
		}

		msg := cmd()
		_, ok := msg.(commands.QuitAddNoteMsg)
		if !ok {
			t.Error("Expected QuitAddNoteMsg from Escape key command")
		}
	})

	t.Run("Enter key with filename creates note and CreateNoteMsg", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		component.textInput.SetValue("mynote")

		keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd := component.ForegroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected ForegroundUpdate to return a command for Enter key with filename")
		}

		msg := cmd()
		createMsg, ok := msg.(commands.CreateNoteMsg)
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
		component := NewComponent(mockRepo)

		component.textInput.SetValue("mynote")

		keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd := component.ForegroundUpdate(keyMsg)

		if cmd == nil {
			t.Fatal("Expected ForegroundUpdate to return a command for Enter key with filename")
		}

		msg := cmd()
		if msg != nil {
			t.Error("Expected nil result when repository fails")
		}
	})

	t.Run("Enter key with empty filename returns nil", func(t *testing.T) {
		mockRepo := &mockRepository{}
		component := NewComponent(mockRepo)

		component.textInput.SetValue("")

		keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd := component.ForegroundUpdate(keyMsg)

		if cmd != nil {
			t.Error("Expected ForegroundUpdate to return nil for Enter key with empty filename")
		}
	})
}
