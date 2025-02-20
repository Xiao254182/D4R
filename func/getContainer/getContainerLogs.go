package getcontainer

import (
	"os/exec"

	"github.com/rivo/tview"
)

func StreamLogs(containerName string, logPanel *tview.TextView, app *tview.Application) {
	cmd := exec.Command("docker", "logs", "-f", "-n", "1000", containerName)

	out, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	if err := cmd.Start(); err != nil {
		return
	}
	defer cmd.Process.Kill()

	buf := make([]byte, 1024)
	for {
		n, err := out.Read(buf)
		if n > 0 {
			app.QueueUpdateDraw(func() {
				logPanel.Write(buf[:n])
				logPanel.ScrollToEnd()
			})
		}
		if err != nil {
			break
		}
	}
}
