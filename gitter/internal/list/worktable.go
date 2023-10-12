package list

import (
	"os"

	tw "github.com/olekukonko/tablewriter"

	"gitter/internal/config"
)

func createTable(layout []config.Column) *Bench {
	names, widths, colors := computeTableParts(layout)
	table := tw.NewWriter(os.Stdout)
	table.SetHeader(names)
	for idx, val := range widths {
		table.SetColMinWidth(idx, val)
	}
	table.SetHeaderColor(colors...)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetBorder(false)
	return &Bench{
		columns: layout,
		tw:      table,
	}
}

type Bench struct {
	columns []config.Column
	tw      *tw.Table
}

func (b *Bench) Append(data map[config.ColumnKind]string, rowColor *RowColor) {
	newRow := make([]string, len(b.columns))
	colors := make([]tw.Colors, len(b.columns))

	for idx, col := range b.columns {
		newRow[idx] = data[col.Kind]
		colors[idx] = rowColor.Get(col.Kind)

		switch col.Wrap {
		case config.Truncate:
			width := col.Width
			if len(newRow[idx]) > width {
				newRow[idx] = newRow[idx][0:width-1] + "…"
			}
		case config.Nothing:
			fallthrough
		default:
		}
	}

	b.tw.Rich(newRow, colors)
}

func (b *Bench) Render() {
	b.tw.Render()
}

func computeTableParts(layout []config.Column) ([]string, []int, []tw.Colors) {
	count := len(layout)
	names := make([]string, 0, count)
	widths := make([]int, 0, count)
	colors := make([]tw.Colors, 0, count)

	for _, col := range layout {
		title := col.Title
		if title == "" {
			title = string(col.Kind)
		}
		names = append(names, title)
		widths = append(widths, col.Width)
		colors = append(colors, tw.Colors{tw.Bold, tw.FgHiBlueColor})
	}

	return names, widths, colors
}

func NewRowColor() *RowColor {
	return &RowColor{
		Colors:       make(map[config.ColumnKind]int),
		Styles:       make(map[config.ColumnKind]int),
		DefaultColor: tw.FgWhiteColor,
		DefaultStyle: tw.Normal,
	}
}

type RowColor struct {
	Colors       map[config.ColumnKind]int
	Styles       map[config.ColumnKind]int
	DefaultColor int
	DefaultStyle int
}

func (c RowColor) Get(kind config.ColumnKind) tw.Colors {
	var style, color int

	if val, ok := c.Styles[kind]; ok {
		style = val
	} else {
		color = c.DefaultStyle
	}
	if val, ok := c.Colors[kind]; ok {
		color = val
	} else {
		color = c.DefaultColor
	}
	return tw.Color(style, color)
}
