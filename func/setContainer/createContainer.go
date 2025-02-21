package setcontainer

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// 创建表单
	form := tview.NewForm().
		AddInputField("Name", "", 30, nil, nil).
		AddInputField("Images", "", 30, nil, nil)

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
	verify := form.HasFocus()
	index := form.GetFormItemByLabel("Images").GetLabel()
	verifystr := strconv.FormatBool(verify)
	form.AddInputField("Images", index, 30, nil, nil)
	form.AddInputField("Images", verifystr, 30, nil, nil)

	if index == "Images" {
		form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch {
			case event.Key() == tcell.KeyTab:
				// 按 Tab 键切换焦点到镜像列表，并添加列表到布局
				flex.AddItem(list, 0, 2, false)
				app.SetFocus(list)
				list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					switch {
					case event.Key() == tcell.KeyEscape:
						// 按 ESC 键切换焦点回表单，隐藏镜像列表
						flex.RemoveItem(list)
						app.SetFocus(form)
						return nil
					case event.Key() == tcell.KeyEnter:
						// 按回车键时，获取选中的镜像并填充到表单的 Images 字段中
						selectImage, _ := list.GetItemText(list.GetCurrentItem())
						// 通过空格拆分字符串
						parts := strings.Fields(selectImage)
						// 拼接镜像名称和标签
						image := fmt.Sprintf("%s:%s", parts[0], parts[1])
						form.RemoveFormItem(1)
						form.AddInputField("Images", image, 30, nil, nil)
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
	}
	// 启动应用
	if err := app.SetRoot(flex, true).Run(); err != nil {
		fmt.Println("启动应用失败:", err)
	}
}

// 获取镜像列表
func getImagesList() []string {
	cmd := exec.Command("docker", "images", "-a")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("执行命令失败:", err)
	}
	// 按行拆分输出并返回镜像列表
	return strings.Split(string(out), "\n")
}
