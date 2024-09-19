package static

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io/ioutil"
	"log"
)

// CreateTextView 从指定路径读取文本文件并返回对应的 TextView
func CreateTextView(filePath string) *tview.TextView {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	textView := tview.NewTextView().
		SetTextAlign(tview.AlignRight).
		SetText(string(content)).
		SetTextColor(tcell.ColorWhite)

	return textView
}
