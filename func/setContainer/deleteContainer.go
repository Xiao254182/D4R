package setcontainer

import (
	getcontainer "D4R/func/getContainer"
	"D4R/types"
	"fmt"
	"os/exec"

	"github.com/rivo/tview"
)

func HandleContainerDeletion(appUI *types.AppUI) {
	index := appUI.ContainerList.GetCurrentItem()
	mainText, _ := appUI.ContainerList.GetItemText(index)
	containerID := getcontainer.ExtractContainerID(mainText)

	if containerID == "" {
		return
	}

	modal := tview.NewModal().
		SetText("是否删除该容器？").
		AddButtons([]string{"取消", "确认删除"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "确认删除" {
				deleteContainer(containerID, appUI)
			}
			appUI.App.SetRoot(appUI.MainPage, true).SetFocus(appUI.ContainerList)
		})

	appUI.App.SetRoot(modal, true)
}
func deleteContainer(containerID string, appUI *types.AppUI) {
	if err := exec.Command("docker", "rm", "-f", containerID).Run(); err != nil {
		showErrorMessage(appUI.App, fmt.Sprintf("删除失败: %v", err))
		return
	}
	refreshContainerList(appUI)
}

func showErrorMessage(app *tview.Application, msg string) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"确定"})
	app.SetRoot(modal, true)
}

func refreshContainerList(appUI *types.AppUI) {
	appUI.ContainerList.Clear()
	containers := getcontainer.GetContainerList()
	for i, name := range containers {
		appUI.ContainerList.AddItem(fmt.Sprintf("%d.%s", i+1, name), "", 0, nil)
	}
}
