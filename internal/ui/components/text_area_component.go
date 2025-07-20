package components

import (
	"elephant/internal/ui/theme"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type TextAreaComponent struct {
	Width, Height int
	textarea      textarea.Model
}

func NewTextArea() *TextAreaComponent {
	ta := textarea.New()
	t := TextAreaComponent{
		Width:    0,
		Height:   0,
		textarea: ta,
	}

	return &t
}

func (t *TextAreaComponent) Init() {
	t.textarea.Focus()
	t.textarea.ShowLineNumbers = false
}

func (t *TextAreaComponent) Update(msg tea.Msg) tea.Cmd {
	newTextarea, cmd := t.textarea.Update(msg)
	t.textarea = newTextarea

	return cmd
}

func (t *TextAreaComponent) View() string {
	textareaView := t.textarea.View()
	return theme.Style.Width(t.Width).Height(t.Height).Render(textareaView)
}

func (t *TextAreaComponent) ResizeWindow(msg tea.WindowSizeMsg) {
	h, v := theme.Style.GetFrameSize()

	t.Width = msg.Width - h
	t.Height = msg.Height - v

	t.textarea.SetWidth(t.Width)
	t.textarea.SetHeight(t.Height)
}

func (t *TextAreaComponent) SetNote(fileContent string) {
	t.textarea.SetValue(fileContent)
}

func (t *TextAreaComponent) GetContent() string {
	return t.textarea.Value()
}
