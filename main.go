package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	headerHeight   = 6
	statsHeight    = 6
	containerWidth = 20
	rightPanel     = 20
)

type AppComponents struct {
	App           *tview.Application
	MainPage      *tview.Flex
	ContainerList *tview.List
	LogPanel      *tview.TextView
	ContainerInfo *tview.TextView
}

type systemInfo struct {
	host, os, arch, dockerVer, cpu, mem string
}

func main() {
	app := tview.NewApplication()
	components := setupUI(app)

	setupGlobalInputHandlers(components)

	if err := app.SetRoot(components.MainPage, true).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}

func setupUI(app *tview.Application) *AppComponents {
	logPanel := createTextViewPanel("Log")
	statsPanel := createTextViewPanel("Stats")
	containerInfoPanel := createTextViewPanel("ContainerInfo")

	containerList := createContainerList(logPanel, statsPanel, containerInfoPanel, app)
	return &AppComponents{
		App:           app,
		MainPage:      createMainLayout(containerList, logPanel, statsPanel, containerInfoPanel),
		ContainerList: containerList,
		LogPanel:      logPanel,
	}
}

func createMainLayout(containerList *tview.List, logPanel, statsPanel, containerInfo *tview.TextView) *tview.Flex {
	header := createHeader()
	outputPanel := createOutputPanel(logPanel)

	separator := tview.NewTextView().SetText(strings.Repeat("- -", 10000)).SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorBlueViolet)

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

func setupGlobalInputHandlers(components *AppComponents) {
	components.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyCtrlC:
			components.App.Stop()
			return nil
		case event.Key() == tcell.KeyCtrlL:
			components.App.SetFocus(components.LogPanel)
			return nil
		case event.Key() == tcell.KeyEscape:
			components.App.SetFocus(components.ContainerList)
			return nil
		case event.Key() == tcell.KeyCtrlI:
			handleContainerExec(components)
			return nil
		case event.Key() == tcell.KeyCtrlD:
			handleContainerDeletion(components)
			return nil
		}
		return event
	})
}

func createHeader() tview.Primitive {
	sysInfo := getSystemInfo()
	titleText := fmt.Sprintf(
		"Host: %s\nOS: %s\nArch: %s\nDocker Version: %s\nCPU: %s\nMem: %s",
		sysInfo.host, sysInfo.os, sysInfo.arch, sysInfo.dockerVer, sysInfo.cpu, sysInfo.mem,
	)

	return tview.NewFlex().
		AddItem(createSystemInfoPanel(titleText), 0, 1, false).
		AddItem(createTipsPanel(), 0, 1, false).
		AddItem(createLogoPanel(), 0, 1, false)
}

func getSystemInfo() systemInfo {
	return systemInfo{
		host:      executeCommand("ip route get 114.114.114.114 | awk '{print $7}'"),
		os:        executeCommand("hostnamectl | grep 'Operating System' | cut -d ':' -f2 | xargs"),
		arch:      executeCommand("hostnamectl | grep 'Architecture' | cut -d ':' -f2 | xargs"),
		dockerVer: executeCommand("docker -v | awk '{print $3}' | sed 's/,//'"),
		cpu:       executeCommand("top -bn1 | grep 'Cpu(s)' | sed 's/.*, *\\([0-9.]*\\)%* id.*/\\1/' | awk '{print 100 - $1 \"%\"}'"),
		mem:       executeCommand("free | awk '/Mem:/ {printf \"%.2f%%\\n\", ($3-$6)/$2 * 100}'"),
	}
}

func executeCommand(cmd string) string {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(out))
}

func createSystemInfoPanel(content string) tview.Primitive {
	return tview.NewTextView().
		SetText(content).
		SetTextColor(tcell.ColorCornflowerBlue).
		SetWrap(false)
}

func createTipsPanel() tview.Primitive {
	return tview.NewTextView().
		SetText(strings.TrimSpace(`
Tips:
  ↑ ↓       切换容器
  Ctrl+C    退出
  Ctrl+L    切换到日志面板
  Ctrl+I    进入容器
  Ctrl+D    删除容器`)).
		SetTextColor(tcell.ColorYellow)
}

func createLogoPanel() tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignRight).
		SetText(strings.TrimSpace(`
	_____  _  _   _____  
   |  __ \| || | |  __ \ 
   | |  | | || |_| |__) |
   | |  | |__   _|  _  / 
   | |__| |  | | | | \ \ 
   |_____/   |_| |_|  \_\
`)).
		SetTextColor(tcell.ColorGreen)
}

func createContainerList(logPanel, statsPanel, containerInfo *tview.TextView, app *tview.Application) *tview.List {
	list := tview.NewList()
	list.SetBorder(true).SetTitle("Containers").SetBorderColor(tcell.ColorLightSkyBlue)

	containers := getContainerList()
	for i, name := range containers {
		list.AddItem(fmt.Sprintf("%d.%s", i+1, name), "", 0, nil)
	}

	var cancelStats context.CancelFunc
	list.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		updateContainerDetails(index, containers, logPanel, statsPanel, app, &cancelStats, containerInfo)
	})

	return list
}

func getContainerList() []string {
	out, err := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}").Output()
	if err != nil {
		return []string{}
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n")
}

func updateContainerDetails(index int, containers []string, logPanel, statsPanel *tview.TextView, app *tview.Application, cancelStats *context.CancelFunc, containerInfo *tview.TextView) {
	logPanel.Clear()
	statsPanel.Clear()

	if index < 0 || index >= len(containers) {
		return
	}
	name := containers[index]

	// 更新容器详情信息
	infoPanel := getContainerInfo(name)
	// 更新右侧容器信息面板
	containerInfo.SetText(infoPanel.GetText(false))

	// 开始日志流
	go streamLogs(name, logPanel, app)

	// 开始统计更新
	if *cancelStats != nil {
		(*cancelStats)()
	}
	ctx, cancel := context.WithCancel(context.Background())
	*cancelStats = cancel
	go updateStats(ctx, name, statsPanel, app)
}

func streamLogs(containerName string, logPanel *tview.TextView, app *tview.Application) {
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

func getContainerInfo(containerName string) *tview.TextView {
	infoPanel := tview.NewTextView()
	infoPanel.SetBorder(true).SetTitle("Container Info").SetBorderColor(tcell.ColorLightSkyBlue)

	// 获取容器的详细信息
	info := getContainerDetails(containerName)
	infoPanel.SetText(info)

	return infoPanel
}

func getContainerDetails(containerName string) string {
	cmd := exec.Command("docker", "inspect", containerName)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("无法获取容器信息: %v", err)
	}

	// 提取我们关心的部分，例如容器的状态、镜像等
	// 这里使用简单的 JSON 解析来提取信息
	// 可以根据需要提取更多字段
	statusCmd := exec.Command("docker", "inspect", "--format", "{{.State.Status}}", containerName)
	status, err := statusCmd.Output()
	if err != nil {
		return fmt.Sprintf("容器状态获取失败: %v", err)
	}

	imageCmd := exec.Command("docker", "inspect", "--format", "{{.Config.Image}}", containerName)
	image, err := imageCmd.Output()
	if err != nil {
		return fmt.Sprintf("容器镜像获取失败: %v", err)
	}

	createdCmd := exec.Command("docker", "inspect", "--format", "{{.Created}}", containerName)
	created, err := createdCmd.Output()
	if err != nil {
		return fmt.Sprintf("容器创建时间获取失败: %v", err)
	}

	// 拼接要显示的信息
	info := fmt.Sprintf("状态: %s\n镜像: %s\n创建时间: %s", string(status), string(image), string(created))
	return info
}

func updateStats(ctx context.Context, containerName string, statsPanel *tview.TextView, app *tview.Application) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			out, err := exec.Command("docker", "stats", "--no-stream", containerName).Output()
			if err == nil {
				app.QueueUpdateDraw(func() {
					statsPanel.SetText(string(out))
				})
			}
		}
	}
}

func handleContainerExec(components *AppComponents) {
	index := components.ContainerList.GetCurrentItem()
	mainText, _ := components.ContainerList.GetItemText(index)
	containerID := extractContainerID(mainText)

	if containerID != "" {
		enterContainer(components.App, containerID)
	}
}

func handleContainerDeletion(components *AppComponents) {
	index := components.ContainerList.GetCurrentItem()
	mainText, _ := components.ContainerList.GetItemText(index)
	containerID := extractContainerID(mainText)

	if containerID == "" {
		return
	}

	modal := tview.NewModal().
		SetText("是否删除该容器？").
		AddButtons([]string{"取消", "确认删除"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "确认删除" {
				deleteContainer(containerID, components)
			}
			components.App.SetRoot(components.MainPage, true).SetFocus(components.ContainerList)
		})

	components.App.SetRoot(modal, true)
}

func extractContainerID(text string) string {
	parts := strings.SplitN(text, ".", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func enterContainer(app *tview.Application, containerID string) {
	cmd := exec.Command("docker", "exec", "-it", containerID, "bash", "-c", "clear; exec /bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	app.Suspend(func() {
		if err := cmd.Run(); err != nil {
			fmt.Printf("执行错误: %v\n", err)
		}
	})
}

func deleteContainer(containerID string, components *AppComponents) {
	if err := exec.Command("docker", "rm", "-f", containerID).Run(); err != nil {
		showErrorMessage(components.App, fmt.Sprintf("删除失败: %v", err))
		return
	}
	refreshContainerList(components)
}

func showErrorMessage(app *tview.Application, msg string) {
	modal := tview.NewModal().
		SetText(msg).
		AddButtons([]string{"确定"})
	app.SetRoot(modal, true)
}

func refreshContainerList(components *AppComponents) {
	components.ContainerList.Clear()
	containers := getContainerList()
	for i, name := range containers {
		components.ContainerList.AddItem(fmt.Sprintf("%d.%s", i+1, name), "", 0, nil)
	}
}

func createTextViewPanel(title string) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	textView.SetTitle(title)
	textView.SetDynamicColors(true)
	textView.SetScrollable(true)
	textView.SetBorderColor(tcell.ColorLightSkyBlue)
	return textView
}

func createOutputPanel(logPanel *tview.TextView) *tview.Flex {
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(logPanel, 0, 1, false)
}
