package main

import (
	env "environment"
	"fmt"
	k8 "k8sinterface"
	"os"
	"view"

	tea "github.com/charmbracelet/bubbletea"
)

func printErr(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func printErrorIfNotNil(err error) {
	if err != nil {
		printErr(err)
	}
}

func main() {
	domen, err := env.ReadFromConfig("app.env", "domen")
	if err != nil {
		q, err := view.NewModelQuestion()
		printErrorIfNotNil(err)
		p := tea.NewProgram(q, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			printErr(err)
		}
	} else {
		clientMicrok8s, err := k8.GetInterfaceProvider(domen)
		printErrorIfNotNil(err)
		w := view.NewModelWaiting(
			func() error {
				return clientMicrok8s.Start()
			})
		p := tea.NewProgram(w, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			printErr(err)
		}
	}
}
