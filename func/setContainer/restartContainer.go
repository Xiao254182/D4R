package setcontainer

import (
	getcontainer "D4R/func/getContainer"
	"D4R/types"
	"fmt"
	"os/exec"

	"github.com/rivo/tview"
)

func HandleContainerRestart(appUI *types.AppUI) {
	index := appUI.ContainerList.GetCurrentItem()
	mainText, _ := appUI.ContainerList.GetItemText(index)
	containerID := getcontainer.ExtractContainerID(mainText)

	if containerID == "" {
		return
	}

	modal := tview.NewModal().
		SetText("是否重启该容器？").
		AddButtons([]string{"取消", "确认重启"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "确认重启" {
				restartContainer(containerID, appUI)
			}
			appUI.App.SetRoot(appUI.MainPage, true).SetFocus(appUI.ContainerList)
		})

	appUI.App.SetRoot(modal, true)
}
func restartContainer(containerID string, appUI *types.AppUI) {
	if err := exec.Command("docker", "restart", containerID).Run(); err != nil {
		showErrorMessage(appUI.App, fmt.Sprintf("重启失败: %v", err))
		return
	}
	refreshContainerList(appUI)
}
