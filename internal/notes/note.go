package notes

type Note struct {
	title, description    string
	filePath, fileContent string
}

func NewNote(title, description string, filePath, fileContent string) Note {
	return Note{
		title:       title,
		description: description,
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
