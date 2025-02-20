package page

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateTextViewPanel(title string) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	textView.SetTitle(title)
	textView.SetDynamicColors(true)
	textView.SetScrollable(true)
	textView.SetBorderColor(tcell.ColorLightSkyBlue)
	return textView
}

func CreateTextViewPanelStats(title string) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	textView.SetTitle(title)
	textView.SetDynamicColors(true)
	textView.SetScrollable(true)
	textView.SetTextAlign(tview.AlignCenter)
	textView.SetBorderColor(tcell.ColorLightSkyBlue)
	return textView
}

func CreateOutputPanel(logPanel *tview.TextView) *tview.Flex {
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(logPanel, 0, 1, false)
}
