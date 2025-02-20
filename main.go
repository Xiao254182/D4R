package main

import (
	"D4R/keyboard"
	"D4R/ui"
	"fmt"
	"os"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	components := ui.SetupLayout(app)

	keyboard.SetupGlobalInputHandlers(components)

	if err := app.SetRoot(components.MainPage, true).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
