package menu

import (
	enter "d4r/exec"
	"d4r/logss"
	"d4r/ps"
	"d4r/rm"
	"github.com/rivo/tview"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// ##############################################docker##########################################
// 处理日志事件
func handleLogEvent(app *tview.Application, table *tview.Table, containers [][]string, logos *tview.TextView, tip *tview.TextView) {
	row, _ := table.GetSelection() // 获取当前选中的行
	if row > 0 {                   // 确保不是表头
		containerID := containers[row-1][0] // 获取选中的容器 ID
		logs, err := logss.GetContainerLogs(containerID)
		if err != nil {
			log.Fatal(err) // 输出错误信息
		}

		logss.ShowLogs(app, logs, func() {
			UpdateTable(app, containers, logos, tip) // 更新表格
		})
	}
}

// 处理进入容器事件
func handleEnterEvent(app *tview.Application, table *tview.Table, containers [][]string, logos *tview.TextView, tip *tview.TextView) {
	row, _ := table.GetSelection()          // 获取当前选中的行
	if row > 0 && row-1 < len(containers) { // 确保不是表头且 row-1 在范围内
		containerID := containers[row-1][0] // 获取选中的容器 ID
		if len(containers[row-1]) > 6 {     // 确保有足够的列
			containerName := containers[row-1][6] // 获取选中的容器名字
			enter.EnterContainer(app, containerID, containerName, func() {
				UpdateTable(app, containers, logos, tip) // 更新表格
			})
		} else {
			log.Print("Error: Container data does not have enough columns.")
		}
	}
}

// 处理从ps界面进入容器事件
func handlePsEnterEvent(app *tview.Application, table *tview.Table, dockerComposePs [][]string, logos *tview.TextView, tip *tview.TextView) {
	row, _ := table.GetSelection() // 获取当前选中的行
	if row > 0 {                   // 确保不是表头
		dockerComposePsName := dockerComposePs[row-1][0] // 获取选中的容器的名字
		containers, _ := ps.GetDockerContainers()
		// 创建一个映射来存储名称到ID的关系
		nameToId := make(map[string]string)
		for _, container := range containers {
			if len(container) > 1 { // 确保有足够的列
				id := container[0]
				name := container[6]
				nameToId[name] = id // 将名称与 ID 关联
			}
		}
		dockerComposePsId, _ := nameToId[dockerComposePsName]

		enter.EnterContainer(app, dockerComposePsId, dockerComposePsName, func() {
			UpdatePsTable(app, dockerComposePs, logos, tip) // 更新表格
		})
	}
}

// 处理删除容器事件
func handleDeleteEvent(app *tview.Application, table *tview.Table, containers [][]string, logos *tview.TextView, tip *tview.TextView) {
	row, _ := table.GetSelection()
	if row > 0 { // 确保不是表头
		containerID := containers[row-1][0] // 获取选中的容器 ID

		// 创建确认框
		modal := tview.NewModal().
			SetText("是否删除该容器？").
			AddButtons([]string{"No", "Yes"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Yes" {
					if err := rm.DeleteContainer(containerID); err != nil {
						log.Fatal(err) // 输出错误信息
					}

					// 删除成功后，获取最新的容器列表
					containers, err := ps.GetDockerContainers() // 获取最新容器信息
					if err != nil {
						log.Fatal(err) // 输出错误信息
					}

					// 刷新表格
					UpdateTable(app, containers, logos, tip) // 调用 updateTable 更新界面
				} else {
					UpdateTable(app, containers, logos, tip) // 更新表格
				}
			})

		app.SetRoot(modal, true) // 显示模态
	}
}

// 处理更新表格信息事件
func UpdateTable(app *tview.Application, containers [][]string, logos *tview.TextView, tip *tview.TextView) {
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(CreateDockerTable(app, containers, logos, tip), 0, 10, true) // 将 logo 传递给 CreateTable

	app.SetRoot(mainFlex, true) // 刷新界面
}

// ##############################################docker-compose##########################################

// 创建 Docker Compose 表格
func handleShowDockerComposePs(app *tview.Application, table *tview.Table, dockerComposePS [][]string, logos *tview.TextView, tip *tview.TextView) *tview.Flex {
	row, _ := table.GetSelection()               // 获取当前选中的行
	if row > 0 && row-1 < len(dockerComposePS) { // 确保不是表头且在范围内
		if len(dockerComposePS[row-1]) > 2 { // 确保有足够的列
			dockerComposePSFile := dockerComposePS[row-1][2] // 获取选中集群的配置文件
			psOutput, err := exec.Command("docker-compose", "-f", dockerComposePSFile, "ps").Output()
			if err != nil {
				log.Fatal(err) // 输出错误信息
			}

			// 解析 docker-compose ps 的输出
			lines := strings.Split(strings.TrimSpace(string(psOutput)), "\n")
			var psOutputs [][]string
			re := regexp.MustCompile(`\s{2,}`)

			for _, line := range lines[1:] { // 跳过表头
				if line = strings.TrimSpace(line); line != "" {
					fields := re.Split(line, -1)
					psOutputs = append(psOutputs, fields)
				}
			}
			return CreateDockerComposePs(app, psOutputs, logos, tip)
		} else {
			log.Print("Error: Docker Compose PS data does not have enough columns.")
		}
	}
	return nil
}

// 处理更新表格信息事件，进入menu后展示
func UpdatePsTable(app *tview.Application, containers [][]string, logos *tview.TextView, tip *tview.TextView) {

	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(CreateDockerComposeTable(app, containers, logos, tip), 0, 10, true) // 将 logo 传递给 CreateTable

	app.SetRoot(mainFlex, true) // 刷新界面
}
