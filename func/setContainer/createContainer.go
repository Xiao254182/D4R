package setcontainer

import (
	"D4R/types"
	"D4R/ui/header"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateContainerFlex 设置容器创建界面的UI
func CreateContainerFlex(appUI *types.AppUI) {
	form := InputContainerForm(appUI)
	header := header.CreateContainerHeader()
	separator := tview.NewTextView().SetText(strings.Repeat("- -", 10000)).SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorLightSkyBlue)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 6, 0, false).
		AddItem(separator, 1, 0, false).
		AddItem(form, 0, 3, true)

	appUI.App.SetRoot(flex, true).SetFocus(form).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			// 处理 ESC 键时的逻辑，停止应用
			SetupGlobalInputHandlers(appUI)
			appUI.App.SetRoot(appUI.MainPage, true).SetFocus(appUI.ContainerList)
			return nil
		}
		return event
	})
}

// InputContainerForm 创建容器表单
func InputContainerForm(appUI *types.AppUI) tview.Primitive {
	app := appUI.App
	var form *tview.Form

	// 创建确认对话框
	popmodal := tview.NewModal().SetFocus(0).
		SetText("是否创建该容器？").
		AddButtons([]string{"取消并返回主页面", "确认"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "确认" {
				createContainer(form)
				refreshContainerList(appUI)
				SetupGlobalInputHandlers(appUI)
				// 手动触发退出键逻辑
				event := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)
				form.InputHandler()(event, func(p tview.Primitive) {})
			}
			app.SetRoot(appUI.MainPage, true).SetFocus(appUI.ContainerList)
		})

	// 创建表单
	form = tview.NewForm().
		AddInputField("Name", "", 30, nil, nil).
		AddInputField("Images", "", 30, nil, nil).
		AddInputField("Port", "", 30, nil, nil).
		AddInputField("Volumes", "", 30, nil, nil).
		AddInputField("Env", "", 30, nil, nil).
		AddInputField("Network", "", 30, nil, nil).
		AddInputField("User", "", 30, nil, nil).
		AddInputField("Workdir", "", 30, nil, nil).
		AddButton("确认创建", func() {
			app.SetRoot(popmodal, true)
		})

	// 创建镜像列表
	list := createImagesList()

	// 创建布局
	flex := tview.NewFlex().AddItem(form, 0, 1, true)

	// 设置焦点切换逻辑
	form.SetInputCapture(handleFormInput(app, appUI, form, list, flex))
	return flex
}

// handleFormInput 处理表单输入逻辑
func handleFormInput(app *tview.Application, appUI *types.AppUI, form *tview.Form, list *tview.List, flex *tview.Flex) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		focusedIndex, _ := form.GetFocusedItemIndex()

		switch event.Key() {
		case tcell.KeyUp:
			// 切换到上一个表单项
			if focusedIndex > 0 {
				app.SetFocus(form.GetFormItem(focusedIndex - 1))
				return nil
			}
		case tcell.KeyDown:
			// 切换到下一个表单项
			if focusedIndex < form.GetFormItemCount()-1 {
				app.SetFocus(form.GetFormItem(focusedIndex + 1))
				return nil
			}
		case tcell.KeyTab:
			// 按 Tab 时展示镜像列表
			if focusedIndex == 1 {
				displayImagesList(app, appUI, form, list, flex)
				return nil
			} else if focusedIndex == 3 {
				displayFileTree(app, form, flex)
				return nil
			}
		}
		return event
	}
}

// displayImagesList 展示镜像列表
func displayImagesList(app *tview.Application, appUI *types.AppUI, form *tview.Form, list *tview.List, flex *tview.Flex) {
	flex.AddItem(list, 0, 2, false)
	app.SetFocus(list)

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			app.SetRoot(appUI.MainPage, true).Run()
			return nil
		case tcell.KeyEnter:
			selectImage, _ := list.GetItemText(list.GetCurrentItem())
			parts := strings.Fields(selectImage)
			if len(parts) >= 2 {
				image := fmt.Sprintf("%s:%s", parts[0], parts[1])
				form.GetFormItem(1).(*tview.InputField).SetText(image)
			}
			flex.RemoveItem(list)
			app.SetFocus(form)
		}
		return event
	})
}

// displayFileTree 展示文件树
func displayFileTree(app *tview.Application, form *tview.Form, flex *tview.Flex) {
	treeView := tview.NewTreeView()
	flex.AddItem(treeView, 0, 2, false)
	app.SetFocus(treeView)

	rootNode := tview.NewTreeNode("/").SetColor(tcell.ColorYellow).SetReference("/").SetExpanded(true)
	treeView.SetRoot(rootNode).SetCurrentNode(rootNode)

	treeView.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}
		path := ref.(string)
		if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {
			if len(node.GetChildren()) == 0 {
				buildTreeNodes(path, node)
			}
			node.SetExpanded(!node.IsExpanded())
		}
	})

	treeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			node := treeView.GetCurrentNode()
			if node == nil {
				return event
			}
			ref := node.GetReference()
			if ref == nil {
				return event
			}
			path := ref.(string)
			form.GetFormItem(3).(*tview.InputField).SetText(path)
			flex.RemoveItem(treeView)
			app.SetFocus(form)
			return nil
		}
		return event
	})
}

// createImagesList 创建镜像列表
func createImagesList() *tview.List {
	list := tview.NewList()
	for _, image := range getImagesList() {
		if image != "" {
			list.AddItem(image, "", 0, nil)
		}
	}
	return list
}

// getImagesList 获取镜像列表
func getImagesList() []string {
	cmd := exec.Command("docker", "images", "-a")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("执行命令失败:", err)
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n")
}

// createContainer 解析表单数据并创建容器
func createContainer(form *tview.Form) {
	name := form.GetFormItem(0).(*tview.InputField).GetText()
	image := form.GetFormItem(1).(*tview.InputField).GetText()
	port := form.GetFormItem(2).(*tview.InputField).GetText()
	volumes := form.GetFormItem(3).(*tview.InputField).GetText()
	env := form.GetFormItem(4).(*tview.InputField).GetText()
	network := form.GetFormItem(5).(*tview.InputField).GetText()
	user := form.GetFormItem(6).(*tview.InputField).GetText()
	workdir := form.GetFormItem(7).(*tview.InputField).GetText()

	if image == "" {
		fmt.Println("错误: 必须指定镜像！")
		return
	}

	args := buildDockerRunArgs(name, port, volumes, env, network, user, workdir, image)
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("创建容器失败: %s\n%s\n", err, strings.TrimSpace(string(output)))
	}
}

// buildDockerRunArgs 构建 docker run 命令参数
func buildDockerRunArgs(name, port, volumes, env, network, user, workdir, image string) []string {
	args := []string{"run", "-d"}
	if name != "" {
		args = append(args, "--name", name)
	}
	if port != "" {
		args = append(args, buildPortArgs(port)...)
	}
	if volumes != "" {
		args = append(args, buildVolumeArgs(volumes)...)
	}
	if env != "" {
		args = append(args, buildEnvArgs(env)...)
	}
	if network != "" {
		args = append(args, "--network", network)
	}
	if user != "" {
		args = append(args, "-u", user)
	}
	if workdir != "" {
		args = append(args, "-w", workdir)
	}
	args = append(args, image)
	return args
}

// buildPortArgs 构建端口参数
func buildPortArgs(port string) []string {
	var args []string
	for _, p := range strings.Split(port, ",") {
		args = append(args, "-p", strings.TrimSpace(p))
	}
	return args
}

// buildVolumeArgs 构建卷参数
func buildVolumeArgs(volumes string) []string {
	var args []string
	for _, v := range strings.Split(volumes, ",") {
		args = append(args, "-v", strings.TrimSpace(v))
	}
	return args
}

// buildEnvArgs 构建环境变量参数
func buildEnvArgs(env string) []string {
	var args []string
	for _, e := range strings.Split(env, ",") {
		args = append(args, "-e", strings.TrimSpace(e))
	}
	return args
}

// buildTreeNodes 递归构建 tview 树节点
func buildTreeNodes(path string, parent *tview.TreeNode) {
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())
		node := tview.NewTreeNode(file.Name()).SetReference(fullPath)

		if file.IsDir() {
			node.SetColor(tcell.ColorYellow).SetExpanded(false)
		} else {
			node.SetColor(tcell.ColorGreen)
		}

		parent.AddChild(node)
	}
}
