package menu

import (
	"d4r/ps"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// 处理 Docker 表格输入事件
func handleDockerTableInputs(event *tcell.EventKey, app *tview.Application, table *tview.Table, containers [][]string, logos *tview.TextView, tip *tview.TextView) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlD:
		dockerCompose, err := ps.GetDockerCompose()
		DockerComposeTable := CreateDockerComposeTable(app, dockerCompose, logos, tip)
		app.SetRoot(DockerComposeTable, true)
		if err != nil {
			return event
		}
	case tcell.KeyRune:
		switch event.Rune() {
		case 'l':
			handleLogEvent(app, table, containers, logos, tip)
		case 'i':
			handleEnterEvent(app, table, containers, logos, tip)
		case 'd':
			handleDeleteEvent(app, table, containers, logos, tip)
		case '\n', '\r':
			return nil
		}
	}
	return event
}

// 处理 Docker Compose 表格输入事件
func handleDockerComposeTableInputs(event *tcell.EventKey, app *tview.Application, table *tview.Table, logos *tview.TextView, tip *tview.TextView) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyESC:
		containers, err := ps.GetDockerContainers()
		if err != nil {
			return event
		}
		DockerTable := CreateDockerTable(app, containers, logos, tip)
		app.SetRoot(DockerTable, true)

	case tcell.KeyRune:
		switch event.Rune() {
		case 'p':
			dockerComposePS, _ := ps.GetDockerCompose()
			if dockerComposePS == nil {
				return nil
			}
			DockerComposeTable := handleShowDockerComposePs(app, table, dockerComposePS, logos, tip)
			app.SetRoot(DockerComposeTable, true)
		case '\n', '\r':
			return nil
		}
	}
	return event
}

// 处理 Docker Compose Ps表格输入事件
func handleDockerComposePsTableInputs(event *tcell.EventKey, app *tview.Application, table *tview.Table, dockerComposePs [][]string, logos *tview.TextView, tip *tview.TextView) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyESC:
		DockerCompose, err := ps.GetDockerCompose()
		if err != nil {
			return event
		}
		DockerComposeTable := CreateDockerComposeTable(app, DockerCompose, logos, tip)
		app.SetRoot(DockerComposeTable, true)

	case tcell.KeyRune:
		switch event.Rune() {
		case 'l':
			handleLogEvent(app, table, dockerComposePs, logos, tip)
		case 'i':
			handlePsEnterEvent(app, table, dockerComposePs, logos, tip)
		case '\n', '\r':
			return nil
		}
	}
	return event
}
