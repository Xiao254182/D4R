package setcontainer

import (
	getcontainer "D4R/func/getContainer"
	"D4R/types"
	"fmt"
	"os"
	"os/exec"

	"github.com/rivo/tview"
)

func HandleContainerExec(appUI *types.AppUI) {
	index := appUI.ContainerList.GetCurrentItem()
	mainText, _ := appUI.ContainerList.GetItemText(index)
	containerID := getcontainer.ExtractContainerID(mainText)

	if containerID != "" {
		enterContainer(appUI.App, containerID)
	}
}
func enterContainer(app *tview.Application, containerID string) {
	cmd := exec.Command("docker", "exec", "-it", containerID, "bash", "-c", "clear; exec /bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	app.Suspend(func() {
		if err := cmd.Run(); err != nil {
			fmt.Printf("执行错误: %v\n", err)
		}
	})
}
