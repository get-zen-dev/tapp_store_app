package shortQuestion

import (
	"constants"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	QUESTION    = "What is your dns?"
	PLACEHOLDER = "salatik.com"

	PART_WINDOW = 3
)

type Question struct {
	question   string
	answer     string
	input      Input
	style      Styles
	dimensions constants.Dimensions
	answered   bool
}

func (q *Question) View() string {
	return lipgloss.Place(
		q.dimensions.Width,
		q.dimensions.Height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			q.question,
			q.style.InputField.Width(q.dimensions.Width/PART_WINDOW).Render(q.input.View()),
		),
	)
}

func (q *Question) SetDimensions(dimensions constants.Dimensions) {
	q.dimensions = dimensions
}

func (q *Question) GetDimensions() constants.Dimensions {
	return q.dimensions
}

func (q *Question) Answered() bool {
	return q.answered
}

func (q *Question) SetAnswered(log bool) {
	q.answered = log
}

func (q *Question) Update(msg tea.Msg) tea.Cmd {
	return q.input.Update(msg)
}

func (q *Question) Input() *Input {
	return &q.input
}

type Styles struct {
	BorderColor lipgloss.AdaptiveColor
	InputField  lipgloss.Style
}

func DefaultStyles() Styles {
	s := Styles{}
	s.BorderColor = lipgloss.AdaptiveColor{Light: "006", Dark: "36"}
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1)
	return s
}

func NewQuestionConcrete() Question {
	return NewQuestion(QUESTION, PLACEHOLDER)
}

func NewQuestion(question, placeholder string) Question {
	q := Question{question: question, answered: false}
	q.input = newInput(placeholder)
	q.style = DefaultStyles()
	q.dimensions = constants.Dimensions{Width: 80, Height: 50}
	return q
}

type Input struct {
	textinput textinput.Model
}

func newInput(placeholder string) Input {
	msg := textinput.New()
	msg.Placeholder = placeholder
	msg.Focus()
	return Input{msg}
}

func (in *Input) Value() string {
	return in.textinput.Value()
}

func (in *Input) Blur() tea.Msg {
	return in.textinput.Blur
}

func (in *Input) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	in.textinput, cmd = in.textinput.Update(msg)
	return cmd
}

func (in *Input) View() string {
	return in.textinput.View()
}
