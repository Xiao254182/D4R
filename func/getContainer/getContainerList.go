package getcontainer

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

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
	apiClient, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	var containerNames []string

	for _, ctr := range containers {
		if len(ctr.Names) > 0 {
			name := strings.TrimPrefix(ctr.Names[0], "/") // 去掉 "/"
			containerNames = append(containerNames, name)
		}
	}
	return containerNames
}

func ExtractContainerID(text string) string {
	parts := strings.SplitN(text, ".", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}
