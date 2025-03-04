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

func InputContainerForm(appUI *types.AppUI) tview.Primitive {
	app := appUI.App

	var form *tview.Form

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
	//创建表单
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
	list := tview.NewList()

	// 获取镜像列表并添加到 list 中
	for _, imageslist := range getImagesList() {
		if imageslist != "" {
			list.AddItem(imageslist, "", 0, nil)
		}
	}

	// 创建布局
	flex := tview.NewFlex().
		AddItem(form, 0, 1, true) // 左侧表单默认显示

	// 设置焦点切换逻辑
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// 获取当前聚焦的表单项索引
		focusedIndex, _ := form.GetFocusedItemIndex()

		switch event.Key() {
		case tcell.KeyUp: // 按 ↑ 切换到上一个表单项
			if focusedIndex > 0 {
				app.SetFocus(form.GetFormItem(focusedIndex - 1))
				return nil
			}
		case tcell.KeyDown: // 按 ↓ 切换到下一个表单项
			if focusedIndex < form.GetFormItemCount()-1 {
				app.SetFocus(form.GetFormItem(focusedIndex + 1))
				return nil
			}
		case tcell.KeyTab: // 仅在第二行按 Tab 时展示列表
			if focusedIndex == 1 {
				flex.AddItem(list, 0, 2, false)
				app.SetFocus(list)

				list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					switch event.Key() {
					case tcell.KeyEscape:
						app.SetRoot(appUI.MainPage, true).Run()
						return nil
					case tcell.KeyEnter:
						// 按回车键时，获取选中的镜像并填充到表单的 Images 字段中
						selectImage, _ := list.GetItemText(list.GetCurrentItem())
						parts := strings.Fields(selectImage)
						if len(parts) >= 2 {
							image := fmt.Sprintf("%s:%s", parts[0], parts[1])
							form.GetFormItem(1).(*tview.InputField).SetText(image)
						}
						// 焦点回到表单
						flex.RemoveItem(list)
						app.SetFocus(form)
					}
					return event
				})
				return nil
			} else {
				if focusedIndex == 3 {
					treeView := tview.NewTreeView()

					flex.AddItem(treeView, 0, 2, false)
					app.SetFocus(treeView)

					// 设置要展示的根目录（当前目录）
					rootPath := "/"
					rootNode := tview.NewTreeNode(rootPath).
						SetColor(tcell.ColorYellow).
						SetReference(rootPath).
						SetExpanded(true)

					treeView.SetRoot(rootNode).SetCurrentNode(rootNode)

					// 处理回车键展开/折叠 或 显示文件路径
					treeView.SetSelectedFunc(func(node *tview.TreeNode) {
						ref := node.GetReference()
						if ref == nil {
							return
						}

						path := ref.(string)

						fileInfo, err := os.Stat(path)
						if err != nil {
							return
						}

						if fileInfo.IsDir() {
							// 如果是目录，展开/折叠
							if len(node.GetChildren()) == 0 {
								buildTreeNodes(path, node) // 加载目录内容
							}
							node.SetExpanded(!node.IsExpanded())
						}
					})

					// 监听 Tab 键，仅在选中文件时填充表单
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
							// 焦点回到表单
							flex.RemoveItem(treeView)
							app.SetFocus(form)
							return nil
						}
						return event
					})
				}
			}
		}
		return event
	})
	return flex
}

// 获取镜像列表
func getImagesList() []string {
	cmd := exec.Command("docker", "images", "-a")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("执行命令失败:", err)
	}
	// 按行拆分输出并返回镜像列表
	return strings.Split(strings.TrimSpace(string(out)), "\n")
}

// 解析表单数据并创建容器
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

	args := []string{"run", "-d"}

	if name != "" {
		args = append(args, "--name", name)
	}
	if port != "" {
		for _, p := range strings.Split(port, ",") {
			args = append(args, "-p", strings.TrimSpace(p))
		}
	}
	if volumes != "" {
		for _, v := range strings.Split(volumes, ",") {
			args = append(args, "-v", strings.TrimSpace(v))
		}
	}
	if env != "" {
		for _, e := range strings.Split(env, ",") {
			args = append(args, "-e", strings.TrimSpace(e))
		}
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
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("创建容器失败: %s\n%s\n", err, strings.TrimSpace(string(output)))
	} else {
	}
}

// 递归构建 tview 树节点
func buildTreeNodes(path string, parent *tview.TreeNode) {
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())
		node := tview.NewTreeNode(file.Name()).SetReference(fullPath)

		if file.IsDir() {
			node.SetColor(tcell.ColorYellow) // 目录黄色
			node.SetExpanded(false)          // 默认不展开
		} else {
			node.SetColor(tcell.ColorGreen) // 文件绿色
		}

		parent.AddChild(node)
	}
}
