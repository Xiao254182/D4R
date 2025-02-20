package getcontainer

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateContainerList(logPanel, statsPanel, containerInfo *tview.TextView, app *tview.Application) *tview.List {
	list := tview.NewList()
	list.SetBorder(true).SetTitle("Containers").SetBorderColor(tcell.ColorLightSkyBlue)

	containers := GetContainerList()
	for i, name := range containers {
		list.AddItem(fmt.Sprintf("%d.%s", i+1, name), "", 0, nil)
	}

	var cancelStats context.CancelFunc
	list.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		UpdateContainerDetails(index, containers, logPanel, statsPanel, app, &cancelStats, containerInfo)
	})

	return list
}

func GetContainerList() []string {
	out, err := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}").Output()
	if err != nil {
		return []string{}
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n")
}

func ExtractContainerID(text string) string {
	parts := strings.SplitN(text, ".", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}
