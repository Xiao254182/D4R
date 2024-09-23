package update

import (
	"d4r/menu"
	"d4r/ps"
	"github.com/rivo/tview"
	"log"
	"time"
)

// 更新容器信息
func UpdateContainers(app *tview.Application, logos, tip *tview.TextView) {
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

			//更新UI
			app.QueueUpdateDraw(func() {
				menu.UpdateTable(app, newContainers, logos, tip)
			})
		}
	}
}

// 更新集群信息
func UpdateDockerComposempose(app *tview.Application, logos, tip *tview.TextView) {
	var previousDockerComposeNames []string

	for {
		time.Sleep(1 * time.Second) // 每秒更新一次

		currentDockerComposeNames, err := ps.GetRunningDockerComposeName() // 获取当前运行的集群
		if err != nil {
			// log.Println("获取集群失败:", err)
			continue
		}

		// 比较当前和之前的容器 ID
		if !equalStrings(previousDockerComposeNames, currentDockerComposeNames) {
			previousDockerComposeNames = currentDockerComposeNames
			newDockerCompose, err := ps.GetDockerCompose()
			if err != nil {
				// log.Println("获取集群失败:", err)
				continue
			}

			// 更新UI
			app.QueueUpdateDraw(func() {
				menu.UpdatePsTable(app, newDockerCompose, logos, tip)
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
