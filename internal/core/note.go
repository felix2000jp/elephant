package core

import (
	"path/filepath"
	"strings"
)

type Note struct {
	title, description    string
	filePath, fileContent string
}

func NewNote(filePath, fileContent string) Note {
	return Note{
		title:       extractTitle(filePath),
		description: extractDescription(fileContent),
		filePath:    filePath,
		fileContent: fileContent,
	}
}

func (n Note) Title() string {
	return n.title
}

func (n Note) Description() string {
	return n.description
}

func (n Note) FilePath() string {
	return n.filePath
}

func (n Note) FileContent() string {
	return n.fileContent
}

func (n Note) FilterValue() string {
	if n.description == "" {
		return n.title
	} else {
		return n.title + " " + n.description
	}
}

func extractTitle(path string) string {
	return strings.TrimSuffix(filepath.Base(path), ".md")
}

func extractDescription(content string) string {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}

	return ""
}
