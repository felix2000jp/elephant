package notes

import (
	"elephant/internal/core"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type EditComponent struct {
	width, height int
	textarea      textarea.Model
	repository    core.Repository

	currentNote core.Note
}

func NewEditComponent(repository core.Repository) EditComponent {
	ta := textarea.New()
	ta.Focus()
	ta.ShowLineNumbers = false

	ec := EditComponent{
		width:      ta.Width(),
		height:     ta.Height(),
		textarea:   ta,
		repository: repository,
	}

	return ec
}

func (ec *EditComponent) Init() tea.Cmd {
	return nil
}

func (ec *EditComponent) BackgroundUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := theme.Style.GetFrameSize()

		ec.width = msg.Width - h
		ec.height = msg.Height - v

		ec.textarea.SetWidth(ec.width)
		ec.textarea.SetHeight(ec.height)

	case ViewNoteMsg:
		ec.currentNote = msg.Note
		ec.textarea.SetValue(msg.Note.FileContent())

	}

	return nil
}

func (ec *EditComponent) ForegroundUpdate(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEsc {
			return func() tea.Msg {
				ec.currentNote = core.NewNote(ec.currentNote.FilePath(), ec.textarea.Value())
				return QuitEditNoteMsg{Note: ec.currentNote}
			}
		}
	}

	var cmd tea.Cmd
	ec.textarea, cmd = ec.textarea.Update(msg)
	return cmd
}

func (ec *EditComponent) View() string {
	listView := ec.textarea.View()
	return theme.Style.Width(ec.width).Height(ec.height).Render(listView)
}
