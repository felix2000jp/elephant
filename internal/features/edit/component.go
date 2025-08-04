package edit

import (
	"elephant/internal/core"
	"elephant/internal/features/commands"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
)

type Component struct {
	width, height int
	textarea      textarea.Model
	repository    core.Repository
	keys          componentKeyMap
	currentNote   core.Note
}

func NewComponent(repository core.Repository) Component {
	keys := newComponentKeyMap()
	ta := textarea.New()
	ta.Focus()
	ta.Prompt = ""
	ta.ShowLineNumbers = false

	ec := Component{
		width:      ta.Width(),
		height:     ta.Height(),
		textarea:   ta,
		repository: repository,
		keys:       keys,
	}

	return ec
}

func (ec *Component) Init() tea.Cmd {
	return nil
}

func (ec *Component) BackgroundUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := theme.Style.GetFrameSize()

		ec.width = msg.Width - h
		ec.height = msg.Height - v

		ec.textarea.SetWidth(ec.width)
		ec.textarea.SetHeight(ec.height)

	case commands.ViewNoteMsg:
		ec.currentNote = msg.Note
		ec.textarea.SetValue(msg.Note.FileContent())

	}

	return nil
}

func (ec *Component) ForegroundUpdate(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, ec.keys.quitEditNote):
			ec.currentNote = core.NewNote(ec.currentNote.FilePath(), ec.textarea.Value())

			return func() tea.Msg {
				err := ec.repository.SaveNote(ec.currentNote)
				if err != nil {
					slog.Error("failed to save note", "error", err)
					return nil
				}

				return commands.QuitEditNoteMsg{Note: ec.currentNote}
			}
		}
	}

	var cmd tea.Cmd
	ec.textarea, cmd = ec.textarea.Update(msg)
	return cmd
}

func (ec *Component) View() string {
	listView := ec.textarea.View()
	return theme.Style.Width(ec.width).Height(ec.height).Render(listView)
}
