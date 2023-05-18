package main

import (
	"fmt"
	"os"
	"requests"

	env "environment"

	tea "github.com/charmbracelet/bubbletea"
	"view"

	k8 "k8sinterface"
)

func main() {
	v8 := k8.GetInterfaceProvider()
	v8.Start()

	env.SetUpEnv()
	list, err := requests.GetListAddons()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	items := view.NewItems()
	for _, v := range list.Value() {
		items.Append(&view.Item{
			Title:          v.Name,
			Status:         "installed",
			CurrentVersion: "1.1.0",
			LastVersion:    "1.1.0"})
	}
	p := tea.NewProgram(view.NewModel(items.GetItems()), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	v8.Stop()
}
