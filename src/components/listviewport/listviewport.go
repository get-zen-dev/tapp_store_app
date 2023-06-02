package listviewport

import (
	"constants"
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"math"
	"style"
)

// Model for controlling the cursor in the table and for displaying the table
type Model struct {
	style           style.Styles
	viewport        viewport.Model
	topBoundId      int
	bottomBoundId   int
	currId          int
	ListItemHeight  int
	NumCurrentItems int
	NumTotalItems   int
}

// Returns new model
func NewModel(styles style.Styles, dimensions constants.Dimensions, numItems, listItemHeight int) Model {
	model := Model{
		style:           styles,
		NumCurrentItems: numItems,
		ListItemHeight:  listItemHeight,
		currId:          0,
		viewport: viewport.Model{
			Width:  dimensions.Width,
			Height: dimensions.Height - style.ListPagerHeight - style.HeaderHeight,
		},
		topBoundId: 0,
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

// Synchronizes the contents of the model for rendering
func (m *Model) SyncViewPort(content string) {
	m.viewport.SetContent(content)
}

func (m *Model) getNumPrsPerPage() int {
	return int(math.Round(float64(float32(m.viewport.Height) / float32(m.ListItemHeight))))
}

func (m *Model) ResetCurrItem() {
	m.currId = 0
	m.viewport.GotoTop()
}

func (m *Model) GetCurrItem() int {
	return m.currId
}

func (m *Model) NextItem() int {
	newId := Min(m.currId+1, m.NumCurrentItems-1)
	newId = Max(newId, 0)
	m.currId = newId

	atBottomOfViewport := m.currId >= m.topBoundId+m.getNumPrsPerPage()
	if atBottomOfViewport {
		m.topBoundId += 1
		m.viewport.LineDown(m.ListItemHeight)
	}
	return m.currId
}

func (m *Model) PrevItem() int {
	m.currId = Max(m.currId-1, 0)

	atTopOfViewport := m.currId < m.topBoundId
	if atTopOfViewport {
		m.topBoundId -= 1
		m.bottomBoundId -= 1
		m.viewport.LineUp(m.ListItemHeight)
	}

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
	m.viewport.Height = dimensions.Height - style.ListPagerHeight - style.HeaderHeight
	m.viewport.Width = dimensions.Width
	m.viewport.SetYOffset(m.currId*m.ListItemHeight - m.getNumPrsPerPage())
}

// Returns a string to be printed to the console
func (m *Model) View() string {
	pagerContent := ""
	if m.NumTotalItems > 0 {
		pagerContent = fmt.Sprintf(
			"%v • %v",
			m.currId+1,
			m.NumTotalItems,
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

// Updates the table style
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
