package menu

import (
	enter "d4r/exec"
	"d4r/logss"
	"d4r/ps"
	"d4r/rm"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
)

// 创建表格
func CreateTable(app *tview.Application, containers [][]string, logos *tview.TextView, tip *tview.TextView) *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	headers := []string{"Container ID", "Image", "Command", "Created", "Status", "Ports", "Names"}
	for i, header := range headers {
		table.SetCell(0, i,
			tview.NewTableCell(header).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter).
				SetExpansion(1)) // 设置可扩展性
	}

	for i, container := range containers {
		for j, field := range container {
			table.SetCell(i+1, j,
				tview.NewTableCell(field).
					SetAlign(tview.AlignCenter).
					SetExpansion(1)) // 设置可扩展性
		}
	}

	// 在 CreateTable 函数中，添加输入捕获
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'l':
			handleLogEvent(app, table, containers, logos, tip) // 用户按下 'l' 键
		case 'i':
			handleEnterEvent(app, table, containers, logos, tip)
		case 'd':
			handleDeleteEvent(app, table, containers, logos, tip) // 用户按下 'd' 键进行删除
		case '\n', '\r':
			return nil // 忽略回车键（换行
		}
		return event // 返回事件以继续处理其他输入
	})

	return table
}

// 处理日志事件
func handleLogEvent(app *tview.Application, table *tview.Table, containers [][]string, logos *tview.TextView, tip *tview.TextView) {
	row, _ := table.GetSelection() // 获取当前选中的行
	if row > 0 {                   // 确保不是表头
		containerID := containers[row-1][0] // 获取选中的容器 ID
		logs, err := logss.GetContainerLogs(containerID)
		if err != nil {
			log.Fatal(err) // 输出错误信息
		}

		logss.ShowLogs(app, logs, func() {
			UpdateTable(app, containers, table, logos, tip) // 更新表格
		})
	}
}

// 处理进入容器事件
func handleEnterEvent(app *tview.Application, table *tview.Table, containers [][]string, logos *tview.TextView, tip *tview.TextView) {
	row, _ := table.GetSelection() // 获取当前选中的行
	if row > 0 {                   // 确保不是表头
		containerID := containers[row-1][0] // 获取选中的容器 ID
		enter.EnterContainer(app, containerID, func() {
			UpdateTable(app, containers, table, logos, tip) // 更新表格
		})
	}
}

// 处理删除容器事件
func handleDeleteEvent(app *tview.Application, table *tview.Table, containers [][]string, logos *tview.TextView, tip *tview.TextView) {
	row, _ := table.GetSelection()
	if row > 0 { // 确保不是表头
		containerID := containers[row-1][0] // 获取选中的容器 ID

		// 创建确认框
		modal := tview.NewModal().
			SetText("是否删除该容器？").
			AddButtons([]string{"No", "Yes"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Yes" {
					if err := rm.DeleteContainer(containerID); err != nil {
						log.Fatal(err) // 输出错误信息
					}

					// 删除成功后，获取最新的容器列表
					containers, err := ps.GetDockerContainers() // 获取最新容器信息
					if err != nil {
						log.Fatal(err) // 输出错误信息
					}

					// 刷新表格
					UpdateTable(app, containers, table, logos, tip) // 调用 updateTable 更新界面
				} else {
					UpdateTable(app, containers, table, logos, tip) // 更新表格
				}
			})

		app.SetRoot(modal, true) // 显示模态
	}
}

// 处理更新表格信息事件
func UpdateTable(app *tview.Application, containers [][]string, table *tview.Table, logos *tview.TextView, tip *tview.TextView) {
	table.Clear() // 清除旧表格内容

	// 更新表格内容
	horizontalFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tip, 0, 1, false).
		AddItem(logos, 0, 2, false)

	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(horizontalFlex, 0, 2, false).
		AddItem(CreateTable(app, containers, logos, tip), 0, 10, true) // 将 logo 传递给 CreateTable

	app.SetRoot(mainFlex, true) // 刷新界面
}
