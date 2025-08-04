package add

import (
	"elephant/internal/core"
	"elephant/internal/features/commands"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
)

type Component struct {
	width, height int
	textInput     textinput.Model
	keys          componentKeyMap
	repository    core.Repository
}

func NewComponent(repository core.Repository) Component {
	keys := newComponentKeyMap()
	ti := textinput.New()
	ti.Placeholder = "Enter note filename (without .md)"
	ti.Focus()

	return Component{
		textInput:  ti,
		keys:       keys,
		repository: repository,
	}
}

func (ac *Component) Init() tea.Cmd {
	return nil
}

func (ac *Component) BackgroundUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := theme.Style.GetFrameSize()
		ac.width = msg.Width - h
		ac.height = msg.Height - v

	case commands.CreateNoteMsg:
		return func() tea.Msg {
			return commands.ViewNoteMsg{Note: msg.Note}
		}
	}

	return nil
}

func (ac *Component) ForegroundUpdate(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, ac.keys.createNote):
			filename := ac.textInput.Value()
			if filename != "" {
				return func() tea.Msg {
					note, err := ac.repository.CreateEmptyNote(filename)
					if err != nil {
						slog.Error("failed to create note", "error", err)
						return nil
					}

					return commands.CreateNoteMsg{Note: note}
				}
			}
		case key.Matches(keyMsg, ac.keys.quitAddNote):
			return func() tea.Msg {
				return commands.QuitAddNoteMsg{}
			}
		}
	}

	var cmd tea.Cmd
	ac.textInput, cmd = ac.textInput.Update(msg)
	return cmd
}

func (ac *Component) View() string {
	content := "Create New Note\n\n" + ac.textInput.View() + "\n\nPress Enter to create, Esc to cancel"
	return theme.Style.Width(ac.width).Height(ac.height).Render(content)
}
