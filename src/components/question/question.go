package question

import (
	"constants"
	"style"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const PART_WINDOW = 3

// Model to hold the question and as well as the input text field
type Question struct {
	question   string
	answer     string
	input      Input
	style      style.StylesQuestion
	dimensions constants.Dimensions
	answered   bool
}

// Returns the question, a text input field that will be printed to 
// the console
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

// Creates a new question model
func NewQuestion(question, placeholder string, style style.StylesQuestion) Question {
	q := Question{question: question, answered: false}
	q.input = newInput(placeholder)
	q.style = style
	q.dimensions = constants.Dimensions{Width: 80, Height: 50}
	return q
}

// Input for text
type Input struct {
	textinput textinput.Model
}

// Creates a new input for text
func newInput(placeholder string) Input {
	msg := textinput.New()
	msg.Placeholder = placeholder
	msg.Focus()
	return Input{msg}
}

// Get input content
func (in *Input) Value() string {
	return in.textinput.Value()
}

func (in *Input) Blur() tea.Msg {
	return in.textinput.Blur
}

// Update content of input
func (in *Input) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	in.textinput, cmd = in.textinput.Update(msg)
	return cmd
}

// Returns a string to be printed to the console
func (in *Input) View() string {
	return in.textinput.View()
}
