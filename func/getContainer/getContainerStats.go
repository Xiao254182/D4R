package getcontainer

import (
	"context"
	"os/exec"
	"time"

	"github.com/rivo/tview"
)

func UpdateStats(ctx context.Context, containerName string, statsPanel *tview.TextView, app *tview.Application) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			out, err := exec.Command("docker", "stats", "--no-stream", containerName, "--format", "table {{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}\t{{.NetIO}}\t{{.BlockIO}}\t{{.PIDs}}").Output()
			if err == nil {
				app.QueueUpdateDraw(func() {
					statsPanel.SetText(string(out))
				})
			}
		}
	}
}
