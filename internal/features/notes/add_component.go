package notes

import (
	"elephant/internal/core"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type addComponent struct {
	width, height int
	textInput     textinput.Model
	repository    core.Repository
}

func newAddComponent(repository core.Repository) addComponent {
	ti := textinput.New()
	ti.Placeholder = "Enter note filename (without .md)"
	ti.Focus()

	return addComponent{
		textInput:  ti,
		repository: repository,
	}
}

func (ac *addComponent) init() tea.Cmd {
	return nil
}

func (ac *addComponent) backgroundUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := theme.Style.GetFrameSize()
		ac.width = msg.Width - h
		ac.height = msg.Height - v

	case CreateNoteMsg:
		note, err := ac.repository.CreateEmptyNote(msg.Filename)
		if err != nil {
			return nil
		}

		return func() tea.Msg {
			return ViewNoteMsg{Note: note}
		}
	}

	return nil
}

func (ac *addComponent) foregroundUpdate(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyEsc:
			return func() tea.Msg {
				return QuitAddNoteMsg{}
			}
		case tea.KeyEnter:
			filename := ac.textInput.Value()
			if filename != "" {
				return func() tea.Msg {
					return CreateNoteMsg{Filename: filename}
				}
			}
		}
	}

	var cmd tea.Cmd
	ac.textInput, cmd = ac.textInput.Update(msg)
	return cmd
}

func (ac *addComponent) view() string {
	content := "Create New Note\n\n" + ac.textInput.View() + "\n\nPress Enter to create, Esc to cancel"
	return theme.Style.Width(ac.width).Height(ac.height).Render(content)
}
