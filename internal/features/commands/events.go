package commands

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

// AddNoteMsg - enter the add note state
type AddNoteMsg struct{}

// QuitAddNoteMsg - quit the add note state
type QuitAddNoteMsg struct{}

// CreateNoteMsg - create a new note with the given filename
type CreateNoteMsg struct{ Note core.Note }
