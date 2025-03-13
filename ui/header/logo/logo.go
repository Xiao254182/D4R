package logo

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateLogoPanel() tview.Primitive {
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
