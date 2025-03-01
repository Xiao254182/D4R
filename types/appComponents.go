package appcomponents

import "github.com/rivo/tview"

type AppComponents struct {
	App           *tview.Application
	MainPage      *tview.Flex
	ContainerList *tview.List
	LogPanel      *tview.TextView
	ContainerInfo *tview.TextView
}
