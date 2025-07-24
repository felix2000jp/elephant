package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNoteRepository(t *testing.T) {
	t.Run("GetAllNotes", func(t *testing.T) {
		tmpDir := createTempDir(t)
		defer removeTempDir(t, tmpDir)

		err := os.WriteFile(filepath.Join(tmpDir, "note1.md"), []byte("# First Note\nThis is the content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err = os.WriteFile(filepath.Join(tmpDir, "note2.md"), []byte("# Second Note\nMore content here"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		service := NewNoteRepository(tmpDir)
		notes, err := service.GetAllNotes()
		if err != nil {
			t.Fatalf("GetAllNotes failed: %v", err)
		}

		if len(notes) != 2 {
			t.Errorf("Expected 2 notes, got %d", len(notes))
		}

		expectedTitles := map[string]string{
			"note1": "First Note",
			"note2": "Second Note",
		}

		for _, note := range notes {
			expectedDesc, exists := expectedTitles[note.Title()]
			if !exists {
				t.Errorf("Unexpected note title: %s", note.Title())
				continue
			}
			if note.Description() != expectedDesc {
				t.Errorf("Expected description '%s', got '%s'", expectedDesc, note.Description())
			}
		}
	})

	t.Run("GetNoteByTitle", func(t *testing.T) {
		tmpDir := createTempDir(t)
		defer removeTempDir(t, tmpDir)

		expectedTitle := "test_note"
		expectedContent := "# Test Note Title\nThis is the test note content"
		testFile := filepath.Join(tmpDir, expectedTitle+".md")

		err := os.WriteFile(testFile, []byte(expectedContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		service := NewNoteRepository(tmpDir)
		note, err := service.GetNoteByTitle(expectedTitle)
		if err != nil {
			t.Fatalf("GetNoteByTitle failed: %v", err)
		}

		if note.Title() != expectedTitle {
			t.Errorf("Expected title '%s', got '%s'", expectedTitle, note.Title())
		}

		if note.FileContent() != expectedContent {
			t.Errorf("Expected content '%s', got '%s'", expectedContent, note.FileContent())
		}

		if note.Description() != "Test Note Title" {
			t.Errorf("Expected description 'Test Note Title', got '%s'", note.Description())
		}

		if note.FilePath() != testFile {
			t.Errorf("Expected file path '%s', got '%s'", testFile, note.FilePath())
		}
	})

	t.Run("SaveNote", func(t *testing.T) {
		tmpDir := createTempDir(t)
		defer removeTempDir(t, tmpDir)

		filePath := filepath.Join(tmpDir, "saved_note.md")
		content := "# Saved Note\nThis is saved content"
		note := NewNote(filePath, content)

		service := NewNoteRepository(tmpDir)
		err := service.SaveNote(note)
		if err != nil {
			t.Fatalf("SaveNote failed: %v", err)
		}

		savedContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read saved file: %v", err)
		}

		if string(savedContent) != content {
			t.Errorf("Expected saved content '%s', got '%s'", content, string(savedContent))
		}
	})
}

func createTempDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "elephant_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	return tmpDir
}

func removeTempDir(t *testing.T, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("Failed to remove temp dir: %v", err)
	}
}
