package main

import (
	"fmt"
	"os"
	"requests"

	env "environment"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"view"
)

func main() {
	env.SetUpEnv()
	list, err := requests.GetListAddons()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	items := view.NewItems([]table.Row{})
	for i, v := range list.Value() {
		if i%2 == 0 {
			items.Append(view.Item{
				Image:          fmt.Sprintf("%s%v", "Image", i),
				Title:          fmt.Sprintf("%s%v", v.Name, i),
				Status:         "installed",
				CurrentVersion: "1.1.0",
				LastVersion:    "1.1.0"})
		} else {
			items.Append(view.Item{
				Image:          fmt.Sprintf("%s%v", "Image", i),
				Title:          fmt.Sprintf("%s%v", v.Name, i),
				Status:         "not installed",
				CurrentVersion: "0.0.8",
				LastVersion:    "0.0.9"})
		}
	}
	p := tea.NewProgram(view.NewModel(items.GetRows()), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
