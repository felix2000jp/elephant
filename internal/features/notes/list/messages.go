package list

import "elephant/internal/core"

type NotesLoadedMsg struct {
	Notes []core.Note
}

type NoteSelectedMsg struct {
	Note core.Note
}
