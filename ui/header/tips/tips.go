package tips

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func MainTipsPanel() tview.Primitive {
	return tview.NewTextView().
		SetText(strings.TrimSpace(`
Tips:
↑ ↓    切换容器	          Ctrl+U 创建一个新的容器
Ctrl+N 切换到容器信息面板 Ctrl+E 进入容器
Ctrl+L 切换到日志面板     Ctrl+D 删除容器
Ctrl+R 重启容器
	`)).
		SetTextColor(tcell.ColorYellow)
}

func CreateContainerTipsPanel() tview.Primitive {
	return tview.NewTextView().
		SetText(strings.TrimSpace(`
Tips:
↑ ↓ 切换表单项	    Tab   提示
Esc 返回主菜单		Enter 下一行/选择   
	`)).
		SetTextColor(tcell.ColorYellow)
}
