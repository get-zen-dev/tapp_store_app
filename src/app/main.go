package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	k8 "k8sinterface"
	"os"
	"view"
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
	clientMicrok8s, err := k8.GetInterfaceProvider("TODO")
	printErrorIfNotNil(err)
	err = clientMicrok8s.Start()
	printErrorIfNotNil(err)
	defer func(clientMicrok8s k8.KuberInterface) {
		printErrorIfNotNil(clientMicrok8s.Stop())
	}(clientMicrok8s)
	m, err := view.NewModel()
	printErrorIfNotNil(err)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		printErr(err)
	}
}
