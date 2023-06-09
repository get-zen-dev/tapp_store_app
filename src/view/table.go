package view

import (
	"constants"
	env "environment"
	"errors"
	he "handleException"
	k8 "k8sinterface"
	"os"
	"requests"
	"style"
	"table"
	"theme"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"

	"log"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ColumnTitleTitle          = "Title"
	ColumnTitleStatus         = "Status"
	ColumnTitleVersion        = "Ver"
	ColumnTitleCurrentVersion = "Current ver"
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

	HeightMessage = 2
)

var (
	headers = []table.Column{
		{Title: ColumnTitleTitle, Width: ColumnMinSizeTitle, MinWidth: ColumnMinSizeTitle, Flex: ColumnFlexTitle},
		{Title: ColumnTitleStatus, Width: ColumnMinSizeStatus, MinWidth: ColumnMinSizeStatus, Flex: ColumnFlexStatus},
		{Title: ColumnTitleVersion, Width: ColumnMinSizeVersion, MinWidth: ColumnMinSizeVersion, Flex: ColumnFlexVersion},
		{Title: ColumnTitleCurrentVersion, Width: ColumnMinSizeCurrentVersion, MinWidth: ColumnMinSizeCurrentVersion, Flex: ColumnFlexCurrentVersion},
		{Title: ColumnTitleDescription, Width: ColumnMinSizeDescription, MinWidth: ColumnMinSizeDescription, Flex: ColumnFlexDescription},
	}

	initStyle = style.InitStyles(*theme.DefaultTheme)

	emptyState = "not found"
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
	table          table.Model
	urls           []string
	style          *style.Styles
	help           help.Model
	lastError      string
	clientMicrok8s k8.KuberInterface
	logger         *log.Logger
	Width          int
	Height         int
}

func NewModelTable(clientMicrok8s k8.KuberInterface, width, height int) (*Model, error) {
	requests.DownloadInfoAddons()
	m := Model{}
	m.clientMicrok8s = clientMicrok8s
	models := env.ReadInfoAddonsModels()
	items := NewItems()
	urls := []string{}
	for _, v := range models.Value() {
		status, curVersion, err := m.getStatusAndCurVersion(v)
		he.PrintErrorIfNotNil(err)
		items.Append(&Item{
			Title:          v.Name,
			Status:         status,
			Version:        v.Version,
			CurrentVersion: curVersion,
			Description:    v.Description})
		url := clientMicrok8s.GetModuleUrl(v.Name)
		urls = append(urls, url.String())
	}
	m.table = table.NewModel(initStyle,
		constants.Dimensions{Width: MinWidth, Height: MinHeight},
		headers,
		items.GetItems(),
		&emptyState)
	m.urls = urls
	m.style = &initStyle
	m.help = help.New()
	m.lastError = ""
	file, _ := os.OpenFile(env.FileLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	m.logger = log.New(file, "\n\n", log.LstdFlags)
	m.Width = width
	m.Height = height
	m.help.Width = width
	m.table.SetDimensions(constants.Dimensions{Width: width, Height: height - constants.Keys.HeightShort - HeightMessage})
	m.table.SyncViewPortContent()
	return &m, nil
}

func (m *Model) getStatusAndCurVersion(v env.Model) (string, string, error) {
	info, err := m.clientMicrok8s.GetCachedModuleInfo(v.Name)
	if err != nil {
		return "", "", err
	}
	curVersion, _ := env.ReadFromConfigCurrentVersion(v.Name)
	if info.IsEnabled {
		if v.Version == curVersion {
			return Installed, curVersion, nil
		} else {
			return Outdated, curVersion, nil
		}
	} else {
		return Deleted, curVersion, nil
	}
}

func (m *Model) updateStatus() error {
	err := m.clientMicrok8s.RefreshInfoCache()
	if err != nil {
		return err
	}
	for i := 0; i < len(m.table.Rows); i++ {
		title := &m.table.Rows[i][Title]
		info, err := m.clientMicrok8s.GetCachedModuleInfo(*title)
		if err != nil {
			return err
		}
		if info.IsEnabled && m.table.Rows[i][Status] == Deleted {
			m.table.Rows[i][Status] = Installed
			version := m.table.Rows[i][Version]
			env.WriteInConfigCurrentVersion(*title, version)
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
	err   error
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.help.Width = msg.Width
		m.table.SetDimensions(constants.Dimensions{Width: msg.Width, Height: msg.Height - constants.Keys.HeightShort - HeightMessage})
		m.table.SyncViewPortContent()
	case Install:
		if msg.err != nil {
			m.logger.Println(msg.err.Error())
			m.lastError = "addon installation error"
		}
		m.table.SyncViewPortContent()
	case Delete:
		if msg.err != nil {
			m.logger.Println(msg.err.Error())
			m.lastError = "addon removal error"
		} else {
			m.table.Rows[msg.index][Status] = Deleted
			env.WriteInConfigCurrentVersion(m.table.Rows[msg.index][Title], "")
			m.table.Rows[msg.index][CurrentVersion] = ""
			m.table.SyncViewPortContent()
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keys.Quit):
			go func() {
				he.PrintErrorIfNotNil(m.clientMicrok8s.Stop())
			}()
			return m, tea.Quit
		case key.Matches(msg, constants.Keys.Up):
			m.table.PrevItem()
		case key.Matches(msg, constants.Keys.Down):
			m.table.NextItem()
		case key.Matches(msg, constants.Keys.Install):
			index := m.table.GetCurrItem()
			name := m.table.Rows[index][Title]
			return m, func() tea.Msg {
				_, errFst := m.clientMicrok8s.InstallModule(name)
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
				return Delete{index, m.clientMicrok8s.RemoveModule(name)}
			}
		case key.Matches(msg, constants.Keys.Update):
			index := m.table.GetCurrItem()
			name := m.table.Rows[index][Title]
			return m, func() tea.Msg {
				if m.table.Rows[index][Version] != m.table.Rows[index][CurrentVersion] {
					m.clientMicrok8s.RemoveModule(name)
					m.clientMicrok8s.InstallModule(name)
					return Install{index, nil}
				}
				return nil
			}
		case key.Matches(msg, constants.Keys.Log):
			content, _ := os.ReadFile(env.FileLog)
			return NewOutputLog(string(content), m, m.Width, m.Height), nil
		}
	}
	return m, nil
}

func (m *Model) View() string {
	message := m.lastError
	m.lastError = ""
	if message == "" {
		message = m.urls[m.table.GetCurrItem()]
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		m.table.View(),
		m.style.Common.FooterStyle.Width(m.help.Width).Render(m.help.View(constants.Keys)),
		m.style.Common.FooterStyle.Width(m.help.Width).Render(m.help.Styles.Ellipsis.Copy().Render(message)),
	)
}
