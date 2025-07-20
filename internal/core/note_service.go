package core

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type NoteService struct {
	basePath string
}

func NewNoteService(basePath string) NoteService {
	return NoteService{basePath: basePath}
}

func (r *NoteService) GetAllNotes() ([]Note, error) {
	pattern := filepath.Join(r.basePath, "*.md")

	files, err := filepath.Glob(pattern)
	if err != nil {
		slog.Error("failed to read markdown files", "error", err)
		return nil, err
	}

	var notes []Note
	for _, filePath := range files {
		title := strings.TrimSuffix(filepath.Base(filePath), ".md")

		content, err := os.ReadFile(filePath)
		if err != nil {
			slog.Warn("failed to read file", "file", filePath, "error", err)
			continue
		}

		fileContent := string(content)
		description := extractDescription(fileContent)
		notes = append(notes, NewNote(title, description, filePath, fileContent))
	}

	slog.Info("loaded notes", "count", len(notes))
	return notes, nil
}

func (r *NoteService) GetNoteByTitle(title string) (Note, error) {
	filePath := filepath.Join(r.basePath, title+".md")

	content, err := os.ReadFile(filePath)
	if err != nil {
		// TODO need to check how to better handle this error
		slog.Error("failed to read note", "file", filePath, "error", err)
		return NewNote(title, "", filePath, ""), err
	}

	fileContent := string(content)
	description := extractDescription(fileContent)
	return NewNote(title, description, filePath, fileContent), nil
}

func (r *NoteService) SaveNote(note Note) error {
	err := os.WriteFile(note.FilePath(), []byte(note.FileContent()), 0644)
	if err != nil {
		slog.Error("failed to save note", "file", note.FilePath, "error", err)
		return err
	}

	return nil
}

func extractDescription(content string) string {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}

	return ""
}
