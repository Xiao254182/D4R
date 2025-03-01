package main

import (
	setcontainer "D4R/func/setContainer"
	"D4R/ui"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/rivo/tview"
)

func main() {
	// 检查是否可以运行 docker 命令
	_, err := exec.Command("docker", "ps").Output()
	if err != nil {
		// 如果执行失败，说明没有安装 Docker
		log.Println("未找到 Docker 环境，D4R已退出")
		os.Exit(1)
	}

	app := tview.NewApplication()
	components := ui.SetupLayout(app)

	setcontainer.SetupGlobalInputHandlers(components)

	if err := app.SetRoot(components.MainPage, true).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
