package notes

import (
	"elephant/internal/core"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type editComponent struct {
	width, height int
	textarea      textarea.Model
	repository    core.Repository

	currentNote core.Note
}

func newEditComponent(repository core.Repository) editComponent {
	ta := textarea.New()
	ta.Focus()
	ta.Prompt = ""
	ta.ShowLineNumbers = false

	ec := editComponent{
		width:      ta.Width(),
		height:     ta.Height(),
		textarea:   ta,
		repository: repository,
	}

	return ec
}

func (ec *editComponent) init() tea.Cmd {
	return nil
}

func (ec *editComponent) backgroundUpdate(msg tea.Msg) tea.Cmd {
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

func (ec *editComponent) foregroundUpdate(msg tea.Msg) tea.Cmd {
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

func (ec *editComponent) view() string {
	listView := ec.textarea.View()
	return theme.Style.Width(ec.width).Height(ec.height).Render(listView)
}
