package constants

import (
	"github.com/charmbracelet/bubbles/key"
)

const ctrlc = "ctrl+c"

// Component for managing keyboard shortcuts
type KeyMap struct {
	Up            key.Binding
	Down          key.Binding
	FirstItem     key.Binding
	LastItem      key.Binding
	TogglePreview key.Binding
	Install       key.Binding
	Delete        key.Binding
	Update        key.Binding
	Help          key.Binding
	Quit          key.Binding
	QuitWithoutQ  key.Binding
	Enter         key.Binding
	HeightShort   int
	HeightFull    int
}

// Returns mini help view that automatically generates itself from your KeyMap
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Install, k.Delete, k.Update}
}

// Returns full help view that automatically generates itself from your KeyMap
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.FirstItem, k.LastItem},
		{k.Help, k.Quit},
	}
}

// Vertical and horizontal dimensions
type Dimensions struct {
	Width  int
	Height int
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	FirstItem: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "first item"),
	),
	LastItem: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "last item"),
	),
	Install: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "install addon"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete addon"),
	),
	Update: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "update addon"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", ctrlc),
		key.WithHelp("q", "quit"),
	),
	QuitWithoutQ: key.NewBinding(
		key.WithKeys(ctrlc),
		key.WithHelp(ctrlc, "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm input"),
	),
	HeightShort: 2,
	HeightFull:  4,
}
