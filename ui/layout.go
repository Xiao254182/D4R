package ui

import (
	getcontainer "D4R/func/getContainer"
	"D4R/types"
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

func SetupLayout(app *tview.Application) *types.AppUI {
	// 统一初始化 UI 组件
	appUI := InitAppUI(app)
	appUI.MainPage = CreateMainLayout(appUI)

	// 获取第一个容器的信息
	if appUI.ContainerList.GetItemCount() > 0 {
		var cancelStats context.CancelFunc
		getcontainer.UpdateContainerDetails(0, getcontainer.GetContainerList(), appUI.LogPanel, appUI.StatsPanel, app, &cancelStats, appUI.ContainerInfo)
	}

	return appUI
}

func InitAppUI(app *tview.Application) *types.AppUI {
	logPanel := page.CreateTextViewPanel("Log")
	statsPanel := page.CreateTextViewPanelStats("Stats")
	containerInfoPanel := page.CreateTextViewPanel("ContainerInfo")

	containerList := getcontainer.CreateContainerList(logPanel, statsPanel, containerInfoPanel, app)

	return &types.AppUI{
		App:           app,
		ContainerList: containerList,
		LogPanel:      logPanel,
		StatsPanel:    statsPanel,
		ContainerInfo: containerInfoPanel,
	}
}

func CreateMainLayout(appUI *types.AppUI) *tview.Flex {
	header := header.CreateHeader()
	separator := tview.NewTextView().SetText(strings.Repeat("- -", 10000)).SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorLightSkyBlue)

	containerDisplay := tview.NewFlex().
		AddItem(appUI.ContainerList, containerWidth, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(appUI.LogPanel, 0, 3, false).
			AddItem(appUI.StatsPanel, statsHeight, 1, false), 0, 2, false).
		AddItem(appUI.ContainerInfo, rightPanel, 1, false)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, headerHeight, 0, false).
		AddItem(separator, 1, 0, false).
		AddItem(containerDisplay, 0, 1, true)
}
