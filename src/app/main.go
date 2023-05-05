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
	k8.CheckInstalledOrInstalKuber()
	env.SetUpEnv()
	list, err := requests.GetListAddons()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(k8.SUStatus())

	items := view.NewItems()
	for i, v := range list.Value() {
		items.Append(&view.Item{
			Number:         fmt.Sprintf("%v", i),
			Title:          fmt.Sprintf("%s%v", v.Name, i),
			Status:         "installed",
			CurrentVersion: "1.1.0",
			LastVersion:    "1.1.0"})
	}
	p := tea.NewProgram(view.NewModel(items.GetItems()), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
