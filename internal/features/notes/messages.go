package notes

import "elephant/internal/core"

// ListNotesMsg - show the list of notes in the base path
type ListNotesMsg struct{ Notes []core.Note }

// ViewNoteMsg - select a single note to view and edit
type ViewNoteMsg struct{ Note core.Note }

// QuitViewNoteMsg - quit the view note state
type QuitViewNoteMsg struct{}

// EditNoteMsg - enter the edit state for the selected note
type EditNoteMsg struct{}

// QuitEditNoteMsg - quit the edit note state
type QuitEditNoteMsg struct{ Note core.Note }
