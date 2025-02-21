package setcontainer

import (
	appcomponents "D4R/func"
	getcontainer "D4R/func/getContainer"
	"fmt"
	"os/exec"

	"github.com/rivo/tview"
)

func HandleContainerRestart(components *appcomponents.AppComponents) {
	index := components.ContainerList.GetCurrentItem()
	mainText, _ := components.ContainerList.GetItemText(index)
	containerID := getcontainer.ExtractContainerID(mainText)

	if containerID == "" {
		return
	}

	modal := tview.NewModal().
		SetText("是否重启该容器？").
		AddButtons([]string{"取消", "确认重启"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "确认重启" {
				restartContainer(containerID, components)
			}
			components.App.SetRoot(components.MainPage, true).SetFocus(components.ContainerList)
		})

	components.App.SetRoot(modal, true)
}
func restartContainer(containerID string, components *appcomponents.AppComponents) {
	if err := exec.Command("docker", "restart", containerID).Run(); err != nil {
		showErrorMessage(components.App, fmt.Sprintf("重启失败: %v", err))
		return
	}
	refreshContainerList(components)
}
