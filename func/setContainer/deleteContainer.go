package setcontainer

import (
	appcomponents "D4R/func"
	getcontainer "D4R/func/getContainer"
	"fmt"
	"os/exec"

	"github.com/rivo/tview"
)

func HandleContainerDeletion(components *appcomponents.AppComponents) {
	index := components.ContainerList.GetCurrentItem()
	mainText, _ := components.ContainerList.GetItemText(index)
	containerID := getcontainer.ExtractContainerID(mainText)

	if containerID == "" {
		return
	}

	modal := tview.NewModal().
		SetText("是否删除该容器？").
		AddButtons([]string{"取消", "确认删除"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "确认删除" {
				deleteContainer(containerID, components)
			}
			components.App.SetRoot(components.MainPage, true).SetFocus(components.ContainerList)
		})

	components.App.SetRoot(modal, true)
}
func deleteContainer(containerID string, components *appcomponents.AppComponents) {
	if err := exec.Command("docker", "rm", "-f", containerID).Run(); err != nil {
		showErrorMessage(components.App, fmt.Sprintf("删除失败: %v", err))
		return
	}
	refreshContainerList(components)
}

func showErrorMessage(app *tview.Application, msg string) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"确定"})
	app.SetRoot(modal, true)
}

func refreshContainerList(components *appcomponents.AppComponents) {
	components.ContainerList.Clear()
	containers := getcontainer.GetContainerList()
	for i, name := range containers {
		components.ContainerList.AddItem(fmt.Sprintf("%d.%s", i+1, name), "", 0, nil)
	}
}
