package listviewport

import (
	"fmt"
	"style"
	"time"

	"common"
	"constants"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	style           style.Styles
	viewport        viewport.Model
	topBoundId      int
	bottomBoundId   int
	currId          int
	ListItemHeight  int
	NumCurrentItems int
	NumTotalItems   int
	LastUpdated     time.Time
	ItemTypeLabel   string
}

func NewModel(style style.Styles, dimensions constants.Dimensions, lastUpdated time.Time, itemTypeLabel string, numItems, listItemHeight int) Model {
	model := Model{
		NumCurrentItems: numItems,
		ListItemHeight:  listItemHeight,
		currId:          0,
		viewport: viewport.Model{
			Width:  dimensions.Width,
			Height: dimensions.Height - common.ListPagerHeight - common.HeaderHeight,
		},
		topBoundId:    0,
		ItemTypeLabel: itemTypeLabel,
		LastUpdated:   lastUpdated,
	}
	model.bottomBoundId = Min(model.NumCurrentItems-1, model.getNumPrsPerPage()-1)
	return model
}

func (m *Model) SetNumItems(numItems int) {
	m.NumCurrentItems = numItems
	m.bottomBoundId = Min(m.NumCurrentItems-1, m.getNumPrsPerPage()-1)
}

func (m *Model) SetTotalItems(total int) {
	m.NumTotalItems = total
}

func (m *Model) SyncViewPort(content string) {
	m.viewport.SetContent(content)
}

func (m *Model) getNumPrsPerPage() int {
	return m.viewport.Height / m.ListItemHeight
}

func (m *Model) ResetCurrItem() {
	m.currId = 0
	m.viewport.GotoTop()
}

func (m *Model) GetCurrItem() int {
	return m.currId
}

func (m *Model) NextItem() int {
	atBottomOfViewport := m.currId >= m.bottomBoundId
	if atBottomOfViewport {
		m.topBoundId += 1
		m.bottomBoundId += 1
		m.viewport.LineDown(m.ListItemHeight)
	}

	newId := Min(m.currId+1, m.NumCurrentItems-1)
	newId = Max(newId, 0)
	m.currId = newId
	return m.currId
}

func (m *Model) PrevItem() int {
	atTopOfViewport := m.currId < m.topBoundId
	if atTopOfViewport {
		m.topBoundId -= 1
		m.bottomBoundId -= 1
		m.viewport.LineUp(m.ListItemHeight)
	}

	m.currId = Max(m.currId-1, 0)
	return m.currId
}

func (m *Model) FirstItem() int {
	m.currId = 0
	m.viewport.GotoTop()
	return m.currId
}

func (m *Model) LastItem() int {
	m.currId = m.NumCurrentItems - 1
	m.viewport.GotoBottom()
	return m.currId
}

func (m *Model) SetDimensions(dimensions constants.Dimensions) {
	m.viewport.Height = dimensions.Height - common.ListPagerHeight - common.HeaderHeight
	m.viewport.Width = dimensions.Width
}

func (m *Model) View() string {
	pagerContent := ""
	if m.NumTotalItems > 0 {
		pagerContent = fmt.Sprintf(
			"%v %v • %v %v/%v • Fetched %v",
			constants.WaitingIcon,
			m.LastUpdated.Format("01/02 15:04:05"),
			m.ItemTypeLabel,
			m.currId+1,
			m.NumTotalItems,
			m.NumCurrentItems,
		)
	}
	viewport := m.viewport.View()
	pager := m.style.ListViewPort.PagerStyle.Copy().Render(pagerContent)
	return lipgloss.NewStyle().
		Width(m.viewport.Width).
		MaxWidth(m.viewport.Width).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			viewport,
			pager,
		))
}

func (m *Model) UpdateStyle(style *style.Styles) {
	m.style = *style
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
