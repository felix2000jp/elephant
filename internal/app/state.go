package app

type State int

const (
	ListState State = iota
	ViewState
	EditState
)
