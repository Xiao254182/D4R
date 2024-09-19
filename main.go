package main

import (
	"d4r/menu"
	"d4r/ps"
	"d4r/static"
	"d4r/update"
	"github.com/rivo/tview"
	"log"
	"os"
	"fmt"
)

// 主函数
func main() {
	app := tview.NewApplication()

	containers, err := ps.GetDockerContainers()
	if err != nil {
		//log.Fatal(err)
		fmt.Println("该系统不存在docker环境或docker服务未启动，请检查docker状态")
		os.Exit(1) // 退出程序，状态码为1表示有错误发生
	}

	logos := static.CreateTextView("./static/logo/logo.txt")
	tip := static.CreateTextView("./static/tip/tip.txt")

	table := menu.CreateTable(app, containers, logos, tip)
	//检测表格的更新
	go update.UpdateContainers(app, table, logos, tip)

	horizontalFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tip, 0, 1, false).
		AddItem(logos, 0, 2, false)

	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(horizontalFlex, 0, 2, false).
		AddItem(table, 0, 10, true)

	if err := app.SetRoot(mainFlex, true).Run(); err != nil {
		log.Fatal(err)
	}
}
