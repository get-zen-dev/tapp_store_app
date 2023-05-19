package main

import (
	env "environment"
	"fmt"
	k8 "k8sinterface"
	"os"
	"requests"

	tea "github.com/charmbracelet/bubbletea"
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
	env.SetUpEnv()
	list, err := requests.GetListAddons()
	if err != nil {
		printErr(err)
	}

	items := view.NewItems()
	for _, v := range list.Value() {
		if v.Name == "common" {
			continue
		}
		info, err := clientMicrok8s.GetModuleInfo(v.Name)
		if err != nil {
			printErr(err)
		}
		status := ""
		if info.IsEnabled {
			status = "✓"
		} else {
			status = "✗"
		}
		items.Append(&view.Item{
			Title:       v.Name,
			Status:      status,
			Description: v.Path})
	}

	p := tea.NewProgram(view.NewModel(items.GetItems()), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		printErr(err)
	}
}
