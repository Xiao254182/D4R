package enter

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// 交互式终端
func EnterContainer(app *tview.Application, containerID string, containerName string, onClose func()) {
	logView := tview.NewTextView().
		SetText(fmt.Sprintf("您当前进入的容器名为:%s,容器id为:%s", containerName, containerID)).
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
			return nil
		}
		return event
	})

	// 设置布局
	modal.SetRect(0, 0, 80, 20)
	logView.SetRect(1, 3, 78, 17)

	// 设置应用程序的根视图为 modal
	app.SetRoot(modal, true).SetFocus(modal)

	// 在 goroutine 中运行 docker exec 命令
	go func() {
		if err := runDockerExec(containerID, logView); err != nil {
			logView.SetText(err.Error())
		}

		// 清空终端并更新 UI
		clearTerminal()
		app.QueueUpdateDraw(func() {
			logView.SetText("") // 清空日志视图
			onClose()           // 调用关闭函数，返回到主界面
		})
	}()

	app.SetRoot(logView, true).SetFocus(logView)
}

// 执行 Docker exec 命令
func runDockerExec(containerID string, logView *tview.TextView) error {
	cmd := exec.Command("docker", "exec", "-it", containerID, "/bin/bash")
	if err := executeCommand(cmd, logView); err != nil {
		// 如果使用 bash 失败，则尝试 sh
		cmd = exec.Command("docker", "exec", "-it", containerID, "sh")
		return executeCommand(cmd, logView)
	}
	return nil
}

// 执行命令并重定向输入输出
func executeCommand(cmd *exec.Cmd, logView *tview.TextView) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// 清空终端
func clearTerminal() {
	clearCmd := exec.Command("clear")
	clearCmd.Stdout = os.Stdout
	clearCmd.Stderr = os.Stderr
	clearCmd.Run()
}
