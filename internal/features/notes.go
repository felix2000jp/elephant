package features

import (
	"elephant/internal/core"
	"elephant/internal/features/add"
	"elephant/internal/features/commands"
	"elephant/internal/features/edit"
	"elephant/internal/features/list"
	"elephant/internal/features/view"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type State int

const (
	ListState State = iota
	ViewState
	EditState
	AddState
)

type NotesFeature struct {
	State         State
	listComponent *list.Component
	viewComponent *view.Component
	editComponent *edit.Component
	addComponent  *add.Component
}

func NewFeature() NotesFeature {
	notesDirectory := getNotesDirectory()

	repository := core.NewNoteRepository(notesDirectory)
	listComponent := list.NewComponent(&repository)
	viewComponent := view.NewComponent(&repository)
	editComponent := edit.NewComponent(&repository)
	addComponent := add.NewComponent(&repository)

	return NotesFeature{
		State:         ListState,
		listComponent: &listComponent,
		viewComponent: &viewComponent,
		editComponent: &editComponent,
		addComponent:  &addComponent,
	}
}

func (nf *NotesFeature) Init() tea.Cmd {
	return tea.Batch(
		nf.listComponent.Init(),
		nf.viewComponent.Init(),
		nf.editComponent.Init(),
		nf.addComponent.Init(),
	)
}

func (nf *NotesFeature) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if _, ok := msg.(commands.ViewNoteMsg); ok {
		nf.State = ViewState
	}
	if _, ok := msg.(commands.QuitViewNoteMsg); ok {
		nf.State = ListState
	}
	if _, ok := msg.(commands.EditNoteMsg); ok {
		nf.State = EditState
	}
	if _, ok := msg.(commands.QuitEditNoteMsg); ok {
		nf.State = ViewState
	}
	if _, ok := msg.(commands.AddNoteMsg); ok {
		nf.State = AddState
	}
	if _, ok := msg.(commands.QuitAddNoteMsg); ok {
		nf.State = ListState
	}

	switch nf.State {
	case ListState:
		cmd = nf.listComponent.ForegroundUpdate(msg)
		cmds = append(cmds, cmd)
	case ViewState:
		cmd = nf.viewComponent.ForegroundUpdate(msg)
		cmds = append(cmds, cmd)
	case EditState:
		cmd = nf.editComponent.ForegroundUpdate(msg)
		cmds = append(cmds, cmd)
	case AddState:
		cmd = nf.addComponent.ForegroundUpdate(msg)
		cmds = append(cmds, cmd)
	}

	cmd = nf.listComponent.BackgroundUpdate(msg)
	cmds = append(cmds, cmd)

	cmd = nf.viewComponent.BackgroundUpdate(msg)
	cmds = append(cmds, cmd)

	cmd = nf.editComponent.BackgroundUpdate(msg)
	cmds = append(cmds, cmd)

	cmd = nf.addComponent.BackgroundUpdate(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (nf *NotesFeature) View() string {
	switch nf.State {
	case ListState:
		return nf.listComponent.View()
	case ViewState:
		return nf.viewComponent.View()
	case EditState:
		return nf.editComponent.View()
	case AddState:
		return nf.addComponent.View()
	default:
		return "Could not render application"
	}
}

func getNotesDirectory() string {
	if dir := os.Getenv("ELEPHANT_NOTES_DIR"); dir != "" {
		return dir
	}
	return ".elephant"
}
