package logss

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os/exec"
)

// 获取特定容器的日志
func GetContainerLogs(containerID string) (string, error) {
	output, err := exec.Command("docker", "logs", containerID).CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// 显示日志的函数
func ShowLogs(app *tview.Application, logs string, onClose func()) {
	logView := tview.NewTextView().
		SetText(logs).
		SetTextColor(tcell.ColorWhite).
		SetScrollable(true)

	modal := tview.NewModal().
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			onClose()
		})

	// 输入捕获，处理 ESC 键
	logView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			onClose()
			return nil // 返回 nil 表示事件已被处理
		}
		return event // 返回原始事件
	})

	// 设置布局
	modal.SetRect(0, 0, 80, 20)
	logView.SetRect(1, 3, 78, 17)

	app.SetRoot(modal, true).SetFocus(modal)
	app.SetRoot(logView, true).SetFocus(logView)
}
