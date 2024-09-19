package update

import (
	"d4r/menu"
	"d4r/ps"
	"github.com/rivo/tview"
	"log"
	"time"
)

// 更新容器信息
func UpdateContainers(app *tview.Application, table *tview.Table, logos, tip *tview.TextView) {
	var previousContainerIDs []string

	for {
		time.Sleep(1 * time.Second) // 每秒更新一次

		currentContainerIDs, err := ps.GetRunningContainerIDs() // 获取当前运行的容器 ID
		if err != nil {
			log.Println("获取容器 ID 失败:", err)
			continue
		}

		// 比较当前和之前的容器 ID
		if !equalStrings(previousContainerIDs, currentContainerIDs) {
			previousContainerIDs = currentContainerIDs
			newContainers, err := ps.GetDockerContainers()
			if err != nil {
				log.Println("获取容器信息失败:", err)
				continue
			}

			// 更新UI
			app.QueueUpdateDraw(func() {
				menu.UpdateTable(app, newContainers, table, logos, tip)
			})
		}
	}
}

// 检查两个字符串切片是否相等
func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
