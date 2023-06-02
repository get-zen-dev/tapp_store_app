package main

import (
	env "environment"
	"fmt"
	he "handleException"
	k8 "k8sinterface"
	"view"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if !k8.CheckIsRootGranted() {
		he.PrintErr(fmt.Errorf("cannot user interface without root privileges"))
	}
	domain, err := env.GetDomain()
	if err != nil {
		q, err := view.NewModelQuestion()
		he.PrintErrorIfNotNil(err)
		p := tea.NewProgram(q, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			he.PrintErr(err)
		}
	} else {
		clientMicrok8s, err := k8.GetInterfaceProvider(domain)
		he.PrintErrorIfNotNil(err)
		w := view.NewModelWaiting(clientMicrok8s, view.KubernetesLaunch)
		p := tea.NewProgram(w, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			he.PrintErr(err)
		}
	}
}
