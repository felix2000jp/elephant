package notes

import (
	"elephant/internal/core"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"log/slog"
)

type ViewComponent struct {
	width, height int
	markdown      viewport.Model
	renderer      *glamour.TermRenderer
	repository    core.Repository

	currentNote core.Note
}

func NewViewComponent(repository core.Repository) ViewComponent {
	renderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(120))
	if err != nil {
		slog.Error("failed to initialize markdown renderer", "error", err)
		panic("failed to initialize markdown renderer")
	}

	vp := viewport.New(0, 0)
	vc := ViewComponent{
		markdown:   vp,
		renderer:   renderer,
		width:      vp.Width,
		height:     vp.Height,
		repository: repository,
	}

	return vc
}

func (vc *ViewComponent) Init() tea.Cmd {
	return nil
}

func (vc *ViewComponent) BackgroundUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := theme.Style.GetFrameSize()

		vc.width = msg.Width - h
		vc.height = msg.Height - v

		vc.markdown.Width = vc.width
		vc.markdown.Height = vc.height

	case ViewNoteMsg:
		vc.currentNote = msg.Note

		content, err := vc.renderer.Render(msg.Note.FileContent())
		if err != nil {
			slog.Error("failed to render markdown", "error", err)
			vc.markdown.SetContent("Could not render content.")
			return nil
		}

		vc.markdown.SetContent(content)

	case QuitEditNoteMsg:
		vc.currentNote = msg.Note

		content, err := vc.renderer.Render(msg.Note.FileContent())
		if err != nil {
			slog.Error("failed to render markdown", "error", err)
			vc.markdown.SetContent("Could not render content.")
			return nil
		}

		vc.markdown.SetContent(content)
	}

	return nil
}

func (vc *ViewComponent) ForegroundUpdate(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEsc {
			return func() tea.Msg {
				return QuitViewNoteMsg{}
			}
		}
		if keyMsg.Type == tea.KeyEnter {
			return func() tea.Msg {
				return EditNoteMsg{}
			}
		}
	}

	var cmd tea.Cmd
	vc.markdown, cmd = vc.markdown.Update(msg)
	return cmd
}

func (vc *ViewComponent) View() string {
	markdownView := vc.markdown.View()
	return theme.Style.Width(vc.width).Height(vc.height).Render(markdownView)
}
