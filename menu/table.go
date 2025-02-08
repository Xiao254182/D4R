package menu

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// 设置表头
func setTableHeaders(table *tview.Table, headers []string) {
	for i, header := range headers {
		table.SetCell(0, i,
			tview.NewTableCell(header).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter).
				SetExpansion(1))
	}
}

// 填充 Docker 容器行
func populateContainerRow(table *tview.Table, rowIndex int, container []string) {
	for j, field := range container {
		table.SetCell(rowIndex, j,
			tview.NewTableCell(field).
				SetAlign(tview.AlignCenter).
				SetExpansion(1))
	}
}

// 填充 Docker Compose 行
func populateComposeRow(table *tview.Table, rowIndex int, container []string) {
	for j, field := range container {
		table.SetCell(rowIndex, j,
			tview.NewTableCell(field).
				SetAlign(tview.AlignCenter).
				SetExpansion(1))
	}
}
