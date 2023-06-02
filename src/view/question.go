package view

import (
	"constants"
	k8 "k8sinterface"
	"question"
	"style"

	env "environment"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type QuestionConcrete struct {
	question question.Question
}

func NewModelQuestion() (*QuestionConcrete, error) {
	m := QuestionConcrete{question.NewQuestion("What is your dns?", "salatik.com", style.InitStylesQuestion())}
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
			domain := m.question.Input().Value()
			env.WriteInConfig("app.yaml", "domain", domain)
			env.GetDomain()
			m.question.SetAnswered(true)
			client, _ := k8.GetInterfaceProvider(domain)
			next := NewModelWaiting(client, KubernetesLaunch)
			next.width = m.question.GetDimensions().Width
			next.height = m.question.GetDimensions().Height
			return next, next.spinner.Tick
		}
	}
	return m, m.question.Update(msg)
}

func (m *QuestionConcrete) View() string {
	return m.question.View()
}
