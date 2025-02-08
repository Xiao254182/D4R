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

func HeaderTitle(widget tview.Primitive, width, height int) tview.Primitive {
	hostCmd := exec.Command("sh", "-c", "ip route get 114.114.114.114 | awk '{print $7}'")
	hostOut, _ := hostCmd.Output()
	host := strings.TrimSpace(string(hostOut))

	osCmd := exec.Command("sh", "-c", "hostnamectl | grep 'Operating System' | cut -d ':' -f2 | xargs")
	osOut, _ := osCmd.Output()
	operatingSystem := strings.TrimSpace(string(osOut))

	archCmd := exec.Command("sh", "-c", "hostnamectl | grep 'Architecture' | cut -d ':' -f2 | xargs")
	archOut, _ := archCmd.Output()
	arch := strings.TrimSpace(string(archOut))

	docker_RevCmd := exec.Command("sh", "-c", "docker -v | awk '{print $3}' | sed 's/,//'")
	docker_RevOut, _ := docker_RevCmd.Output()
	docker_ver := strings.TrimSpace(string(docker_RevOut))

	cpuCmd := exec.Command("sh", "-c", "top -bn1 | grep 'Cpu(s)' | sed 's/.*, *\\([0-9.]*\\)%* id.*/\\1/' | awk '{print 100 - $1 \"%\"}'")
	cpuOut, _ := cpuCmd.Output()
	cpu := strings.TrimSpace(string(cpuOut))

	memCmd := exec.Command("sh", "-c", "free | awk '/Mem:/ {printf \"%.2f%%\\n\", ($3-$6)/$2 * 100}'")
	memOut, _ := memCmd.Output()
	mem := strings.TrimSpace(string(memOut))

	titleText := fmt.Sprintf("tag: 占位\nHost: %s\nOS: %s\narch: %s\nD4R_vev: v2.0\ndocker_ver: %s\nCPU: %s\nMem: %s", host, operatingSystem, arch, docker_ver, cpu, mem)
	lineCount := strings.Count(titleText, "\n") + 3
	if height < lineCount {
		height = lineCount // 动态调整高度
	}

	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				AddItem(tview.NewTextView().SetText(titleText), 0, 1, false).
				AddItem(widget, width, 0, true),
			height, 0, true).
		AddItem(nil, 0, 1, false)
}

func main() {
	app := tview.NewApplication()

	imagelist := createImageList()
	logPanel := createTextViewPanel(app, "Log")
	statsPanel := createTextViewPanel(app, "Stats")
	containerList := createApplication(imagelist, logPanel, statsPanel, app)
	outputPanel := createOutputPanel(logPanel)
	headerTitle := HeaderTitle(tview.NewTextView(), 0, 1)

	MainPage := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerTitle, 7, 0, false).
		AddItem(tview.NewFlex().
			AddItem(containerList, 20, 1, true).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(outputPanel, 0, 3, false).
				AddItem(statsPanel, 6, 1, false), 0, 2, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 20, 1, false), 0, 1, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
			return nil
		case tcell.KeyCtrlL:
			app.SetFocus(logPanel)
			return nil
		case tcell.KeyEscape:
			app.SetFocus(containerList)
			return nil
		case tcell.KeyCtrlI:
			index := containerList.GetCurrentItem()
			mainText, _ := containerList.GetItemText(index)

			if mainText != "" {
				parts := strings.SplitN(mainText, ".", 2)
				if len(parts) == 2 {
					mainText = parts[1]
				}
				fmt.Println("Executing docker exec for container:", mainText)

				EnterContainer(app, strings.TrimSpace(mainText), mainText, MainPage, containerList, func() {
					app.SetRoot(MainPage, true).SetFocus(containerList)
				})
			}
			return nil
		case tcell.KeyCtrlD:
			index := containerList.GetCurrentItem()
			mainText, _ := containerList.GetItemText(index)

			modal := tview.NewModal().
				SetText("是否删除该容器？").
				AddButtons([]string{"No", "Yes"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Yes" {
						if mainText != "" {
							parts := strings.SplitN(mainText, ".", 2)
							if len(parts) == 2 {
								mainText = parts[1]
							}
							DeleteContainer(app, strings.TrimSpace(mainText), mainText, MainPage, containerList, func() {
								app.SetRoot(MainPage, true).SetFocus(containerList)
							})
						} else {
							fmt.Println("No container selected")
						}
					} else {
						app.SetRoot(MainPage, true).SetFocus(containerList)
					}
				})
			app.SetRoot(modal, true) // 显示模态
		}
		return event
	})
	// 运行 TUI
	if err := app.SetRoot(MainPage, true).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}

func createContainerList() (commandList *tview.List) {
	commandList = tview.NewList()
	commandList.SetBorder(true).SetTitle("Containers")
	commandList.ShowSecondaryText(true)
	return commandList
}

func createImageList() (commandList *tview.List) {
	commandList = tview.NewList()
	commandList.SetBorder(true).SetTitle("Images")
	commandList.ShowSecondaryText(false)
	commandList.SetSelectedFocusOnly(true)
	return commandList
}

func createApplication(imageList *tview.List, logPanel *tview.TextView, statsPanel *tview.TextView, app *tview.Application) *tview.List {
	containerList := createContainerList()

	dockernameout, err := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}").Output()
	if err == nil {
		containernames := strings.Split(strings.TrimSpace(string(dockernameout)), "\n")
		for i, containerName := range containernames {
			name := containerName
			containerList.AddItem(
				fmt.Sprintf("%d.%s", i+1, name),
				"",
				rune(0),
				nil,
			)
		}

		var cancelStats context.CancelFunc

		containerList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			imageList.Clear()
			logPanel.Clear()
			statsPanel.Clear()

			name := containernames[index]

			cmd := exec.Command("docker", "inspect", "--format", "{{.Config.Image}}", name)
			out, err := cmd.Output()
			if err == nil {
				imageName := strings.TrimSpace(string(out))
				imageList.AddItem(imageName, "", 0, nil)
			}

			go func() {
				logCmd := exec.Command("docker", "logs", "-f", "-n", "1000", name)
				logOut, err := logCmd.StdoutPipe()
				if err != nil {
					return
				}
				if err := logCmd.Start(); err != nil {
					return
				}
				buf := make([]byte, 1024)
				for {
					n, err := logOut.Read(buf)
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
			}()

			if cancelStats != nil {
				cancelStats()
			}

			ctx, cancel := context.WithCancel(context.Background())
			cancelStats = cancel

			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						statsCmd := exec.Command("docker", "stats", "--no-stream", name)
						statsOut, err := statsCmd.Output()
						if err == nil {
							app.QueueUpdateDraw(func() {
								statsPanel.SetText(string(statsOut))
							})
						}
						time.Sleep(1 * time.Second)
					}
				}
			}()
		})
	}

	return containerList
}

func createTextViewPanel(app *tview.Application, name string) (panel *tview.TextView) {
	panel = tview.NewTextView()
	panel.SetBorder(true).SetTitle(name)
	panel.SetChangedFunc(func() {
		app.Draw()
	})
	panel.SetDynamicColors(true)
	panel.SetScrollable(true)
	return panel
}

func createOutputPanel(logPanel *tview.TextView) (outputPanel *tview.Flex) {
	outputPanel = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(logPanel, 0, 1, false)
	return outputPanel
}

func EnterContainer(app *tview.Application, containerID string, containerName string, flex *tview.Flex, containerList *tview.List, callback func()) {
	shellView := tview.NewTextView()

	go func() {
		if err := runDockerExec(app, containerID, shellView); err != nil {
			app.QueueUpdateDraw(func() {
				shellView.SetText(fmt.Sprintf("Error: %s", err.Error()))
			})
		}
		callback()
	}()
}

func runDockerExec(app *tview.Application, containerID string, logView *tview.TextView) error {
	// 先执行 clear 命令
	clearCmd := exec.Command("clear")
	clearCmd.Stdout = os.Stdout
	clearCmd.Run()

	// 然后执行 docker exec 命令
	cmd := exec.Command("docker", "exec", "-it", containerID, "bash", "-c", "clear; exec /bin/bash")
	if err := executeCommand(app, cmd, logView); err != nil {
		cmd = exec.Command("docker", "exec", "-it", containerID, "bash", "-c", "clear; exec /bin/bash")
		return executeCommand(app, cmd, logView)
	}
	return nil
}

func executeCommand(app *tview.Application, cmd *exec.Cmd, shellView *tview.TextView) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	app.QueueUpdateDraw(func() {
		modal := tview.NewFlex().
			AddItem(shellView, 0, 1, true).
			AddItem(tview.NewBox().
				SetBorder(false), 10, 1, false)
		modal.Clear()
		app.SetRoot(modal, true).SetFocus(shellView)
	})

	err := cmd.Run()
	return err
}

func DeleteContainer(app *tview.Application, containerID string, containerName string, flex *tview.Flex, containerList *tview.List, callback func()) {
	cmd := exec.Command("docker", "rm", "-f", containerID)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error deleting container:", err)
	} else {
		app.SetRoot(flex, true).SetFocus(containerList)
	}
}
