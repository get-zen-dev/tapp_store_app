package table

import (
	"constants"
	"listviewport"
	"style"

	"github.com/charmbracelet/lipgloss"
)

// Column in a table
type Column struct {
	Title    string
	Width    int
	MinWidth int
	Flex     int
}

// Row to table
type Row []string

// Table Model
type Model struct {
	style         style.Styles
	Columns       []Column
	Rows          []Row
	EmptyState    *string
	dimensions    constants.Dimensions
	minDimensions constants.Dimensions
	rowsViewport  listviewport.Model
}

// Create a new table model
func NewModel(style style.Styles, dimensions constants.Dimensions, columns []Column, rows []Row, emptyState *string) Model {
	return Model{
		style:         style,
		Columns:       columns,
		Rows:          rows,
		EmptyState:    emptyState,
		minDimensions: dimensions,
		rowsViewport:  listviewport.NewModel(style, dimensions, len(rows), 2),
	}
}

// Returns a string to print to the console
func (m Model) View() string {
	header := m.renderHeader()
	body := m.renderBody()
	return lipgloss.JoinVertical(lipgloss.Left, header, body)
}

// Set dimensions for the table
func (m *Model) SetDimensions(dimensions constants.Dimensions) {
	m.dimensions = constants.Dimensions{Width: listviewport.Max(m.minDimensions.Width, dimensions.Width),
		Height: listviewport.Max(m.minDimensions.Height, dimensions.Height)}
	m.rowsViewport.SetDimensions(constants.Dimensions{
		Width:  m.dimensions.Width,
		Height: m.dimensions.Height,
	})
}

func (m *Model) ResetCurrItem() {
	m.rowsViewport.ResetCurrItem()
}

// Returns the index of the current row
func (m *Model) GetCurrItem() int {
	return m.rowsViewport.GetCurrItem()
}

// Move index to previous row
func (m *Model) PrevItem() int {
	currItem := m.rowsViewport.PrevItem()
	m.SyncViewPortContent()

	return currItem
}

// Move index to next row
func (m *Model) NextItem() int {
	currItem := m.rowsViewport.NextItem()
	m.SyncViewPortContent()

	return currItem
}

func (m *Model) FirstItem() int {
	currItem := m.rowsViewport.FirstItem()
	m.SyncViewPortContent()

	return currItem
}

func (m *Model) LastItem() int {
	currItem := m.rowsViewport.LastItem()
	m.SyncViewPortContent()

	return currItem
}

// Synchronizes the contents of the model for listviewport
func (m *Model) SyncViewPortContent() {
	headerColumns := m.renderHeaderColumns()
	renderedRows := make([]string, 0, len(m.Rows))
	for i := range m.Rows {
		renderedRows = append(renderedRows, m.renderRow(i, headerColumns))
	}

	m.rowsViewport.SyncViewPort(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	)
}

func (m *Model) SetRows(rows []Row) {
	m.Rows = rows
	m.rowsViewport.SetNumItems(len(m.Rows))
	m.SyncViewPortContent()
}

func (m *Model) OnLineDown() {
	m.rowsViewport.NextItem()
}

func (m *Model) OnLineUp() {
	m.rowsViewport.PrevItem()
}

func (m *Model) getShownColumns() []Column {
	shownColumns := make([]Column, 0, len(m.Columns))
	return append(shownColumns, m.Columns...)
}

func (m *Model) renderHeaderColumns() []string {
	shownColumns := m.getShownColumns()
	renderedColumns := make([]string, len(shownColumns))
	allFlex := 0
	for _, column := range shownColumns {
		allFlex += column.Flex
	}
	leftoverWidth := m.dimensions.Width
	widthFlex := leftoverWidth / allFlex
	width := 0
	for i, column := range shownColumns {
		if i != len(renderedColumns)-1 {
			width = widthFlex * column.Flex
		} else {
			width = leftoverWidth
		}
		column.Width = width
		leftoverWidth -= width
		renderedColumns[i] = m.style.Table.TitleCellStyle.Copy().
			Width(column.Width).
			MaxWidth(column.Width).
			Render(column.Title)
	}
	return renderedColumns
}

func (m *Model) renderHeader() string {
	headerColumns := m.renderHeaderColumns()
	header := lipgloss.JoinHorizontal(lipgloss.Top, headerColumns...)
	return m.style.Table.HeaderStyle.Copy().
		Width(m.dimensions.Width).
		MaxWidth(m.dimensions.Width).
		Height(style.TableHeaderHeight).
		MaxHeight(style.TableHeaderHeight).
		Render(header)
}

func (m *Model) renderBody() string {
	bodyStyle := lipgloss.NewStyle().
		Height(m.dimensions.Height).
		MaxWidth(m.dimensions.Width)

	if len(m.Rows) == 0 && m.EmptyState != nil {
		return bodyStyle.Render(*m.EmptyState)
	}

	return m.rowsViewport.View()
}

func (m *Model) renderRow(rowId int, headerColumns []string) string {
	var style lipgloss.Style

	if m.rowsViewport.GetCurrItem() == rowId {
		style = m.style.Table.SelectedCellStyle
	} else {
		style = m.style.Table.CellStyle
	}

	renderedColumns := make([]string, 0, len(m.Columns))

	headerColId := 0

	for i := range m.Columns {
		colWidth := lipgloss.Width(headerColumns[headerColId])
		content := m.Rows[rowId][i]
		if colWidth < len(content) + 1 {
			content = content[0 : listviewport.Max(colWidth-5, 0)]
			content += "..."
		}
		renderedCol := style.Copy().Width(colWidth).MaxWidth(colWidth).Height(1).MaxHeight(1).Render(content)
		renderedColumns = append(renderedColumns, renderedCol)
		headerColId++
	}

	return m.style.Table.RowStyle.Copy().
		MaxWidth(m.dimensions.Width).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedColumns...))
}

func (m *Model) UpdateStyle(style *style.Styles) {
	m.style = *style
	m.rowsViewport.UpdateStyle(style)
}

func (m *Model) UpdateTotalItemsCount(count int) {
	m.rowsViewport.SetTotalItems(count)
}
