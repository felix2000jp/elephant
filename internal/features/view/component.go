package view

import (
	"elephant/internal/core"
	"elephant/internal/features/commands"
	"elephant/internal/theme"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"log/slog"
)

type Component struct {
	width, height int
	markdown      viewport.Model
	renderer      *glamour.TermRenderer
	keyMap        customKeyMap
	repository    core.Repository

	currentNote core.Note
}

func NewComponent(repository core.Repository) Component {
	keyMap := newCustomKeyMap()
	renderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(120))
	if err != nil {
		slog.Error("failed to initialize markdown renderer", "error", err)
		panic("failed to initialize markdown renderer")
	}

	vp := viewport.New(0, 0)
	vp.KeyMap = keyMap.baseKeyMap

	vc := Component{
		width:      vp.Width,
		height:     vp.Height,
		markdown:   vp,
		renderer:   renderer,
		keyMap:     keyMap,
		repository: repository,
	}

	return vc
}

func (vc *Component) Init() tea.Cmd {
	return nil
}

func (vc *Component) BackgroundUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := theme.Style.GetFrameSize()

		vc.width = msg.Width - h
		vc.height = msg.Height - v

		vc.markdown.Width = vc.width
		vc.markdown.Height = vc.height

	case commands.ViewNoteMsg:
		vc.currentNote = msg.Note

		content, err := vc.renderer.Render(msg.Note.FileContent())
		if err != nil {
			slog.Error("failed to render markdown", "error", err)
			vc.markdown.SetContent("Could not render content.")
			return nil
		}

		vc.markdown.SetContent(content)

	case commands.QuitEditNoteMsg:
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

func (vc *Component) ForegroundUpdate(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, vc.keyMap.quitViewNote):
			return func() tea.Msg {
				return commands.QuitViewNoteMsg{}
			}
		case key.Matches(keyMsg, vc.keyMap.editNote):
			return func() tea.Msg {
				return commands.EditNoteMsg{}
			}
		}
	}

	var cmd tea.Cmd
	vc.markdown, cmd = vc.markdown.Update(msg)
	return cmd
}

func (vc *Component) View() string {
	markdownView := vc.markdown.View()
	return theme.Style.Width(vc.width).Height(vc.height).Render(markdownView)
}
