package tips

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateTipsPanel() tview.Primitive {
	return tview.NewTextView().
		SetText(strings.TrimSpace(`
Tips:
↑ ↓    切换容器	          Ctrl+U 创建一个新的容器
Ctrl+N 切换到容器信息面板 Ctrl+I 进入容器
Ctrl+L 切换到日志面板     Ctrl+D 删除容器
	`)).
		SetTextColor(tcell.ColorYellow)
}
