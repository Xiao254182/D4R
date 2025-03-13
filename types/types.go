package types

import "github.com/rivo/tview"

type AppUI struct {
	App           *tview.Application
	MainPage      *tview.Flex
	ContainerList *tview.List
	LogPanel      *tview.TextView
	StatsPanel    *tview.TextView
	ContainerInfo *tview.TextView
}
