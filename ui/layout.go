package ui

import (
	appcomponents "D4R/func"
	getcontainer "D4R/func/getContainer"
	"D4R/ui/header"
	"D4R/ui/page"
	"context"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	headerHeight   = 6
	statsHeight    = 6
	containerWidth = 15
	rightPanel     = 35
)

func SetupLayout(app *tview.Application) *appcomponents.AppComponents {
	logPanel := page.CreateTextViewPanel("Log")
	statsPanel := page.CreateTextViewPanelStats("Stats")
	containerInfoPanel := page.CreateTextViewPanel("ContainerInfo")

	containerList := getcontainer.CreateContainerList(logPanel, statsPanel, containerInfoPanel, app)
	components := &appcomponents.AppComponents{
		App:           app,
		MainPage:      createMainLayout(containerList, logPanel, statsPanel, containerInfoPanel),
		ContainerList: containerList,
		LogPanel:      logPanel,
		ContainerInfo: containerInfoPanel,
	}

	// 获取第一个容器的信息
	if containerList.GetItemCount() > 0 {
		var cancelStats context.CancelFunc
		getcontainer.UpdateContainerDetails(0, getcontainer.GetContainerList(), logPanel, statsPanel, app, &cancelStats, containerInfoPanel)
	}

	return components
}

func createMainLayout(containerList *tview.List, logPanel, statsPanel, containerInfo *tview.TextView) *tview.Flex {
	header := header.CreateHeader()
	outputPanel := page.CreateOutputPanel(logPanel)

	separator := tview.NewTextView().SetText(strings.Repeat("- -", 10000)).SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorLightSkyBlue)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, headerHeight, 0, false).
		AddItem(separator, 1, 0, false).
		AddItem(tview.NewFlex().
			AddItem(containerList, containerWidth, 1, true).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(outputPanel, 0, 3, false).
				AddItem(statsPanel, statsHeight, 1, false), 0, 2, false).
			AddItem(containerInfo, rightPanel, 1, false), 0, 1, true)
}
