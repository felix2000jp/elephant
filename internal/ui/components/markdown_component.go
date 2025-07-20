package components

import (
	"elephant/internal/core"
	"elephant/internal/ui/theme"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"log/slog"
)

type MarkdownComponent struct {
	Width, Height int
	markdown      viewport.Model
	renderer      *glamour.TermRenderer
}

func NewMarkdownComponent() *MarkdownComponent {
	vp := viewport.New(0, 0)
	m := MarkdownComponent{
		markdown: vp,
		Width:    vp.Width,
		Height:   vp.Height,
	}

	return &m
}

func (m *MarkdownComponent) Init() {
	renderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(120))
	if err != nil {
		panic("failed to initialize markdown renderer")
	}

	m.renderer = renderer
}

func (m *MarkdownComponent) Update(msg tea.Msg) tea.Cmd {
	newNoteMarkdown, cmd := m.markdown.Update(msg)
	m.markdown = newNoteMarkdown

	return cmd
}

func (m *MarkdownComponent) View() string {
	markdownView := m.markdown.View()
	return theme.Style.Width(m.Width).Height(m.Height).Render(markdownView)
}

func (m *MarkdownComponent) ResizeWindow(msg tea.WindowSizeMsg) {
	h, v := theme.Style.GetFrameSize()

	m.Width = msg.Width - h
	m.Height = msg.Height - v

	m.markdown.Width = m.Width
	m.markdown.Height = m.Height
}

func (m *MarkdownComponent) SetContent(note core.Note) {
	content, err := m.renderer.Render(note.FileContent())
	if err != nil {
		slog.Error("failed to render markdown", "error", err)
		m.markdown.SetContent("Could not render markdown.")
		return
	}

	m.markdown.SetContent(content)
}
