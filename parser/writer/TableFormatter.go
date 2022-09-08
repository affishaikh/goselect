package writer

import (
	"fmt"
	"github.com/bndr/gotabulate"
	"goselect/parser/executor"
	"goselect/parser/projection"
)

type TableFormatter struct{}

func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

func (tableFormatter TableFormatter) Format(projections *projection.Projections, rows *executor.EvaluatingRows) string {
	var displayableRows [][]string

	iterator := rows.RowIterator()
	for iterator.HasNext() {
		var attributes []string
		for _, attribute := range iterator.Next().AllAttributes() {
			attributes = append(attributes, attribute.GetAsString())
		}
		displayableRows = append(displayableRows, attributes)
	}

	return tableFormatter.renderContentTable(projections.DisplayableAttributes(), displayableRows) +
		tableFormatter.renderFooterTable(displayableRows)
}

func (tableFormatter TableFormatter) renderContentTable(attributes []string, displayableRows [][]string) string {
	table := gotabulate.Create(displayableRows)
	table.SetHeaders(attributes)
	table.SetMaxCellSize(45)
	table.SetWrapStrings(true)
	table.SetAlign("left")

	return table.Render("grid")
}

func (tableFormatter TableFormatter) renderFooterTable(displayableRows [][]string) string {
	table := gotabulate.Create([]string{fmt.Sprintf("Total Rows: %v", len(displayableRows))})
	table.SetHeaders([]string{""})
	table.SetAlign("left")
	return table.Render("grid")
}
