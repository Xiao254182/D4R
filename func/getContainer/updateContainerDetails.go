package getcontainer

import (
	"context"

	"github.com/rivo/tview"
)

func UpdateContainerDetails(index int, containers []string, logPanel, statsPanel *tview.TextView, app *tview.Application, cancelStats *context.CancelFunc, containerInfo *tview.TextView) {
	logPanel.Clear()
	statsPanel.Clear()

	if index < 0 || index >= len(containers) {
		return
	}
	name := containers[index]

	// 更新容器详情信息
	infoPanel := CreateContainerOut(name)
	containerInfo.SetText(infoPanel.GetText(false))

	go StreamLogs(name, logPanel, app)

	if *cancelStats != nil {
		(*cancelStats)()
	}
	ctx, cancel := context.WithCancel(context.Background())
	*cancelStats = cancel
	go UpdateStats(ctx, name, statsPanel, app)
}
