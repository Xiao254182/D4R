package setcontainer

import (
	appcomponents "D4R/func"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InputContainerForm(components *appcomponents.AppComponents) tview.Primitive {
	app := components.App

	var form *tview.Form

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
			createContainer(form)
		}).
		AddButton("取消", func() {
			app.SetRoot(form, true)
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

		if event.Key() == tcell.KeyTab && focusedIndex == 1 {
			fmt.Println("test")
			// 仅在第二行按 Tab 时展示列表
			// 按 Tab 键切换焦点到镜像列表，并添加列表到布局
			flex.AddItem(list, 0, 2, false)
			app.SetFocus(list)

			list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				switch event.Key() {
				case tcell.KeyEscape:
					// 按 ESC 键切换焦点回表单，隐藏镜像列表
					flex.RemoveItem(list)
					app.SetFocus(form)
					return nil
				case tcell.KeyEnter:
					// 按回车键时，获取选中的镜像并填充到表单的 Images 字段中
					selectImage, _ := list.GetItemText(list.GetCurrentItem())
					// 通过空格拆分字符串
					parts := strings.Fields(selectImage)
					if len(parts) >= 2 { // 确保有足够的字段
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
		fmt.Println("容器创建成功:", strings.TrimSpace(string(output)))
	}
}
