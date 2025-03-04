package header

import (
	"D4R/ui/header/logo"
	"D4R/ui/header/systemInfo"
	"D4R/ui/header/tips"
	"fmt"

	"github.com/rivo/tview"
)

func createHeader(tipsPanel tview.Primitive) tview.Primitive {
	sysInfo := systeminfo.GetSystemInfo()
	titleText := fmt.Sprintf(
		"Host: %s\nOS: %s\nArch: %s\nDocker Version: %s\nCPU: %s\nMem: %s",
		sysInfo.Host, sysInfo.OS, sysInfo.Arch, sysInfo.DockerVer, sysInfo.CPU, sysInfo.Mem,
	)

	return tview.NewFlex().
		AddItem(systeminfo.CreateSystemInfoPanel(titleText), 0, 1, false).
		AddItem(tipsPanel, 0, 1, false).
		AddItem(logo.CreateLogoPanel(), 0, 1, false)
}

func MainHeader() tview.Primitive {
	return createHeader(tips.MainTipsPanel())
}

func CreateContainerHeader() tview.Primitive {
	return createHeader(tips.CreateContainerTipsPanel())
}
