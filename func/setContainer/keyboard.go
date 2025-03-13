package setcontainer

import (
	"D4R/types"

	"github.com/gdamore/tcell/v2"
)

func SetupGlobalInputHandlers(appUI *types.AppUI) {
	appUI.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyCtrlC:
			appUI.App.Stop()
			return nil
		case event.Key() == tcell.KeyCtrlL:
			appUI.App.SetFocus(appUI.LogPanel)
			return nil
		case event.Key() == tcell.KeyEscape:
			appUI.App.SetFocus(appUI.ContainerList)
			return nil
		case event.Key() == tcell.KeyCtrlE:
			HandleContainerExec(appUI)
			return nil
		case event.Key() == tcell.KeyCtrlD:
			HandleContainerDeletion(appUI)
			return nil
		case event.Key() == tcell.KeyCtrlN:
			appUI.App.SetFocus(appUI.ContainerInfo)
			return nil
		case event.Key() == tcell.KeyCtrlR:
			HandleContainerRestart(appUI)
			return nil
		case event.Key() == tcell.KeyCtrlU:
			CreateContainerFlex(appUI)
			return nil
		}
		return event
	})
}
