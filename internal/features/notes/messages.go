package notes

// TODO think about if this file is needed or the messages should be defined next to the sender

type ViewNoteMsg struct {
	NoteTitle string
}

type EditNoteMsg struct {
	NoteTitle string
}
