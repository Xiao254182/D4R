package setcontainer

import (
	appcomponents "D4R/func"
	"fmt"
	"os"

	"github.com/rivo/tview"
)

func createContainerInfo(components *appcomponents.AppComponents) *tview.Application {
	form := tview.NewForm().
		AddInputField("Name", "", 20, nil, nil).
		AddInputField("Image", "", 20, nil, nil).
		AddInputField("Port Mapping", "", 20, nil, nil).
		AddInputField("Network", "", 20, nil, nil).
		AddInputField("Volume", "", 20, nil, nil).
		AddInputField("Environment", "", 20, nil, nil).
		AddInputField("User", "", 20, nil, nil).
		AddInputField("Working Dir", "", 20, nil, nil)
	form.AddButton("OK", func() {
		//todo
	}).
		AddButton("Cancel", func() {
			components.App.SetRoot(components.MainPage, true).SetFocus(components.ContainerList)
		})
	app := tview.NewApplication()
	if err := app.SetRoot(form, true).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
	return app
}
