package view

import (
	"constants"
	env "environment"
	"errors"
	k8 "k8sinterface"
	"requests"
	"style"
	"table"
	"theme"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ColumnTitleTitle          = "Title"
	ColumnTitleStatus         = "Status"
	ColumnTitleVersion        = "Version"
	ColumnTitleCurrentVersion = "Now"
	ColumnTitleDescription    = "Description"

	ColumnFlexTitle          = 6
	ColumnFlexStatus         = 4
	ColumnFlexVersion        = 4
	ColumnFlexCurrentVersion = 4
	ColumnFlexDescription    = 6

	ColumnMinSizeTitle          = 8
	ColumnMinSizeStatus         = 8
	ColumnMinSizeVersion        = 8
	ColumnMinSizeCurrentVersion = 8
	ColumnMinSizeDescription    = 10

	Title          = 0
	Status         = 1
	Version        = 2
	CurrentVersion = 3
	Description    = 4

	MinWidth  = 50
	MinHeight = 10

	Installed = "✓"
	Deleted   = "✗"
	Outdated  = "—"

	HeightLastLink = 2
)

var (
	headers = []table.Column{
		{Title: ColumnTitleTitle, Width: ColumnMinSizeTitle, MinWidth: ColumnMinSizeTitle, Flex: ColumnFlexTitle},
		{Title: ColumnTitleStatus, Width: ColumnMinSizeStatus, MinWidth: ColumnMinSizeStatus, Flex: ColumnFlexStatus},
		{Title: ColumnTitleVersion, Width: ColumnMinSizeVersion, MinWidth: ColumnMinSizeVersion, Flex: ColumnFlexVersion},
		{Title: ColumnTitleCurrentVersion, Width: ColumnMinSizeCurrentVersion, MinWidth: ColumnMinSizeCurrentVersion, Flex: ColumnFlexCurrentVersion},
		{Title: ColumnTitleDescription, Width: ColumnMinSizeDescription, MinWidth: ColumnMinSizeDescription, Flex: ColumnFlexDescription},
	}

	domen, _          = env.ReadFromConfig("app.env", "domen")
	clientMicrok8s, _ = k8.GetInterfaceProvider(domen)

	currentVersion = "current_version.yaml"
)

type Item struct {
	Title          string
	Status         string
	Version        string
	CurrentVersion string
	Description    string
}

type Items struct {
	items []table.Row
}

func NewItems() *Items {
	return &Items{}
}

func (i *Items) Append(item *Item) {
	i.items = append(i.items, makeRow(item))
}

func makeRow(item *Item) []string {
	return []string{
		item.Title,
		item.Status,
		item.Version,
		item.CurrentVersion,
		item.Description,
	}
}

func (i *Items) GetItems() []table.Row {
	return i.items
}

type Model struct {
	table    table.Model
	style    *style.Styles
	help     help.Model
	lastLink string
}

func NewModel() (*Model, error) {
	err := requests.DownloadInfoAddons()
	if err != nil {
		return nil, err
	}
	models := env.ReadInfoAddonsModels()
	items := NewItems()
	for _, v := range models.Value() {
		info, err := clientMicrok8s.GetCachedModuleInfo(v.Name)
		if err != nil {
			return nil, err
		}
		status := ""
		curVersion, _ := env.ReadFromConfig(currentVersion, v.Name)
		if info.IsEnabled {
			if v.Version == curVersion {
				status = Installed
			} else {
				status = Outdated
			}
		} else {
			status = Deleted
		}
		items.Append(&Item{
			Title:          v.Name,
			Status:         status,
			Version:        v.Version,
			CurrentVersion: curVersion,
			Description:    v.Description})
	}
	s := style.InitStyles(*theme.DefaultTheme)
	emptyState := "not found"
	m := Model{
		table: table.NewModel(s,
			constants.Dimensions{Width: MinWidth, Height: MinHeight},
			headers,
			items.GetItems(),
			&emptyState),
		style:    &s,
		help:     help.New(),
		lastLink: "",
	}
	return &m, nil
}

func (m *Model) updateStatus() error {
	err := clientMicrok8s.RefreshInfoCache()
	if err != nil {
		return err
	}
	for i := 0; i < len(m.table.Rows); i++ {
		title := m.table.Rows[i][Title]
		info, err := clientMicrok8s.GetCachedModuleInfo(title)
		if err != nil {
			return err
		}
		curVersion, _ := env.ReadFromConfig(currentVersion, title)
		if info.IsEnabled && curVersion == "" {
			m.table.Rows[i][Status] = Installed
			version := m.table.Rows[i][Version]
			env.WriteInConfig(currentVersion, title, version)
			m.table.Rows[i][CurrentVersion] = version
		}
	}
	return nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

type Install struct {
	index int
	err   error
}

type Delete struct {
	index int
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.table.SetDimensions(constants.Dimensions{Width: msg.Width, Height: msg.Height - constants.Keys.HeightShort - HeightLastLink})
		m.table.SyncViewPortContent()
	case Install:
		if msg.err != nil {
			m.lastLink = "error: " + msg.err.Error()
		} else {
			m.lastLink = "link: https://codeforces.com/" + m.table.Rows[msg.index][Title]
		}
		m.table.SyncViewPortContent()
	case Delete:
		m.table.Rows[msg.index][Status] = Deleted
		env.WriteInConfig(currentVersion, m.table.Rows[msg.index][Title], "")
		m.table.Rows[msg.index][CurrentVersion] = ""
		m.lastLink = ""
		m.table.SyncViewPortContent()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, constants.Keys.Up):
			m.table.PrevItem()
		case key.Matches(msg, constants.Keys.Down):
			m.table.NextItem()
		case key.Matches(msg, constants.Keys.Install):
			index := m.table.GetCurrItem()
			name := m.table.Rows[index][Title]
			return m, func() tea.Msg {
				_, errFst := clientMicrok8s.InstallModule(name)
				errSnd := m.updateStatus()
				var err error = nil
				if errFst != nil || errSnd != nil {
					err = errors.Join(errFst, errSnd)
				}
				return Install{index, err}
			}
		case key.Matches(msg, constants.Keys.Delete):
			index := m.table.GetCurrItem()
			name := m.table.Rows[index][Title]
			return m, func() tea.Msg {
				clientMicrok8s.RemoveModule(name)
				return Delete{index}
			}
		case key.Matches(msg, constants.Keys.Update):
			index := m.table.GetCurrItem()
			name := m.table.Rows[index][Title]
			return m, func() tea.Msg {
				if m.table.Rows[index][Version] != m.table.Rows[index][CurrentVersion] {
					clientMicrok8s.RemoveModule(name)
					clientMicrok8s.InstallModule(name)
					return Install{index, nil}
				}
				return nil
			}
		}
	}
	return m, nil
}

func (m *Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		m.table.View(),
		m.style.Common.FooterStyle.Width(m.help.Width).Render(m.help.View(constants.Keys)),
		m.style.Common.FooterStyle.Width(m.help.Width).Render(m.help.Styles.Ellipsis.Copy().Render(m.lastLink)),
	)
}
