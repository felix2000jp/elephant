package core

import (
	"log/slog"
	"os"
	"path/filepath"
)

type Repository interface {
	GetAllNotes() ([]Note, error)
	GetNoteByTitle(title string) (Note, error)
	SaveNote(note Note) error
	CreateEmptyNote(filename string) (Note, error)
}

type NoteRepository struct {
	basePath string
}

func NewNoteRepository(basePath string) NoteRepository {
	return NoteRepository{basePath: basePath}
}

func (r *NoteRepository) GetAllNotes() ([]Note, error) {
	pattern := filepath.Join(r.basePath, "*.md")

	files, err := filepath.Glob(pattern)
	if err != nil {
		slog.Error("failed to read markdown files", "error", err)
		return nil, err
	}

	var notes []Note
	for _, filePath := range files {
		content, err := os.ReadFile(filePath)
		if err != nil {
			slog.Warn("failed to read file", "file", filePath, "error", err)
			continue
		}

		fileContent := string(content)
		notes = append(notes, NewNote(filePath, fileContent))
	}

	slog.Info("loaded notes", "count", len(notes))
	return notes, nil
}

func (r *NoteRepository) GetNoteByTitle(title string) (Note, error) {
	fileName := title + ".md"
	filePath := filepath.Join(r.basePath, fileName)

	content, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("failed to read note by title", "title", title, "file", filePath, "error", err)
		return Note{}, err
	}

	return NewNote(filePath, string(content)), nil
}

func (r *NoteRepository) SaveNote(note Note) error {
	err := os.WriteFile(note.FilePath(), []byte(note.FileContent()), 0644)
	if err != nil {
		slog.Error("failed to save note", "file", note.FilePath(), "error", err)
		return err
	}

	return nil
}

func (r *NoteRepository) CreateEmptyNote(filename string) (Note, error) {
	if filepath.Ext(filename) != ".md" {
		filename = filename + ".md"
	}

	filePath := filepath.Join(r.basePath, filename)

	err := os.WriteFile(filePath, []byte(""), 0644)
	if err != nil {
		slog.Error("failed to create empty note", "file", filePath, "error", err)
		return Note{}, err
	}

	return NewNote(filePath, ""), nil
}
