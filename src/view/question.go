package view

import (
	"constants"
	"shortQuestion"

	env "environment"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type QuestionConcrete struct {
	question shortQuestion.Question
}

func NewQuestion() (*QuestionConcrete, error) {
	m := QuestionConcrete{shortQuestion.NewQuestionConcreteDomen()}
	return &m, nil
}

func (m QuestionConcrete) Init() tea.Cmd {
	return nil
}

func (m *QuestionConcrete) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.question.Answered() {
		return m, tea.Quit
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.question.SetDimensions(constants.Dimensions{Width: msg.Width, Height: msg.Height})
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keys.QuitWithoutQ):
			return m, tea.Quit
		case key.Matches(msg, constants.Keys.Enter):
			domen := m.question.Input().Value()
			env.WriteInConfig("domen", domen)
			m.question.SetAnswered(true)
			return m, m.question.Input().Blur
		}
	}
	return m, m.question.Update(msg)
}

func (m *QuestionConcrete) View() string {
	return m.question.View()
}
