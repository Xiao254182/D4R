package getcontainer

import (
	"bufio"
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/rivo/tview"
)

func StreamLogs(containerName string, logPanel *tview.TextView, app *tview.Application) {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation()) // 根据您的Docker版本选择合适的API版本
	if err != nil {
		return
	}

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "1000",
	}

	out, err := cli.ContainerLogs(context.Background(), containerName, options)
	if err != nil {
		return
	}

	defer out.Close()

	// 从日志流中读取并处理日志行
	// 将日志流转换为扫描器
	scanner := bufio.NewScanner(out)
	//不断读取
	for scanner.Scan() {
		line := scanner.Bytes()

		// 将[]byte切片转换为字符串输出
		// 因为docker前8个字符都是不可见字符，所以需要去掉，从第九个开始读取
		resultString := string(line)[8:]
		if resultString != "" {
			app.QueueUpdateDraw(func() {
				logPanel.Write([]byte(resultString + "\n"))
				logPanel.ScrollToEnd()
			})
		}

	}
}
