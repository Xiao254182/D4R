package menu

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// 创建 Docker 表格
func CreateDockerTable(app *tview.Application, containers [][]string, logos *tview.TextView, tip *tview.TextView) *tview.Flex {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	headers := []string{"Container ID", "Image", "Command", "Created", "Status", "Ports", "Names"}
	setTableHeaders(table, headers)

	for i, container := range containers {
		populateContainerRow(table, i+1, container)
	}

	// 输入捕获
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return handleDockerTableInputs(event, app, table, containers, logos, tip)
	})

	// 创建一个水平 Flex 布局，将 tip 和 logos 组合在一起
	horizontalFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tip, 0, 1, false).  // 将 tip 添加到左边
		AddItem(logos, 0, 2, false) // 将 logos 添加到右边

	// 创建主 Flex 布局
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(horizontalFlex, 0, 3, false). // 将水平布局添加到上面
		AddItem(table, 0, 10, true)           // 将表格添加到剩余空间

	return mainFlex

}

// 创建 Docker Compose 表格
func CreateDockerComposeTable(app *tview.Application, containers [][]string, logos *tview.TextView, tip *tview.TextView) *tview.Flex {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	headers := []string{"NAME", "STATUS", "CONFIG FILES"}
	setTableHeaders(table, headers)

	for i, container := range containers {
		populateComposeRow(table, i+1, container)
	}

	// 输入捕获
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return handleDockerComposeTableInputs(event, app, table, logos, tip)
	})

	// 创建一个水平 Flex 布局，将 tip 和 logos 组合在一起
	horizontalFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tip, 0, 1, false).  // 将 tip 添加到左边
		AddItem(logos, 0, 2, false) // 将 logos 添加到右边

	// 创建主 Flex 布局
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(horizontalFlex, 0, 3, false). // 将水平布局添加到上面
		AddItem(table, 0, 10, true)           // 将表格添加到剩余空间

	return mainFlex
}

// 创建 Docker Compose Ps表格
func CreateDockerComposePs(app *tview.Application, dockerComposePs [][]string, logos *tview.TextView, tip *tview.TextView) *tview.Flex {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	headers := []string{"NAME", "IMAGE", "COMMAND", "SERVICE", "CREATED", "STATUS", "PORTS"}
	setTableHeaders(table, headers)

	for i, dockerComposePs := range dockerComposePs {
		populateComposeRow(table, i+1, dockerComposePs)
	}

	// 输入捕获
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return handleDockerComposePsTableInputs(event, app, table, dockerComposePs, logos, tip)
	})

	// 创建一个水平 Flex 布局，将 tip 和 logos 组合在一起
	horizontalFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tip, 0, 1, false).  // 将 tip 添加到左边
		AddItem(logos, 0, 2, false) // 将 logos 添加到右边

	// 创建主 Flex 布局
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(horizontalFlex, 0, 3, false). // 将水平布局添加到上面
		AddItem(table, 0, 10, true)           // 将表格添加到剩余空间

	return mainFlex

}
