package keyboard

import (
	setcontainer "D4R/func/setContainer"
	appcomponents "D4R/types"

	"github.com/gdamore/tcell/v2"
)

func SetupGlobalInputHandlers(components *appcomponents.AppComponents) {
	components.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyCtrlC:
			components.App.Stop()
			return nil
		case event.Key() == tcell.KeyCtrlL:
			components.App.SetFocus(components.LogPanel)
			return nil
		case event.Key() == tcell.KeyEscape:
			components.App.SetFocus(components.ContainerList)
			return nil
		case event.Key() == tcell.KeyCtrlE:
			setcontainer.HandleContainerExec(components)
			return nil
		case event.Key() == tcell.KeyCtrlD:
			setcontainer.HandleContainerDeletion(components)
			return nil
		case event.Key() == tcell.KeyCtrlN:
			components.App.SetFocus(components.ContainerInfo)
			return nil
		case event.Key() == tcell.KeyCtrlR:
			setcontainer.HandleContainerRestart(components)
			return nil
		case event.Key() == tcell.KeyCtrlU:
			setcontainer.CreateContainerFlex(components)
			return nil
		}
		return event
	})
}
