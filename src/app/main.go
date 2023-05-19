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

func main() {
	clientMicrok8s := k8.GetInterfaceProvider()
	err := clientMicrok8s.Start()
	if err != nil {
		printErr(err)
	}
	/*	defer func(clientMicrok8s k8.KuberInterface) {
		err := clientMicrok8s.Stop()
		if err != nil {
			printErr(err)
		}
	}(clientMicrok8s)*/
	m, err := view.NewModel()
	if err != nil {
		printErr(err)
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		printErr(err)
	}
}
