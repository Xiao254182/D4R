package systeminfo

import (
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SystemInfo struct {
	Host, OS, Arch, DockerVer, CPU, Mem string
}

func CreateSystemInfoPanel(content string) tview.Primitive {
	return tview.NewTextView().
		SetText(content).
		SetTextColor(tcell.ColorCornflowerBlue).
		SetWrap(false)
}

func executeCommand(cmd string) string {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(out))
}

func GetSystemInfo() SystemInfo {
	return SystemInfo{
		Host:      executeCommand("ip route get 114.114.114.114 | awk '{print $7}'"),
		OS:        executeCommand("hostnamectl | grep 'Operating System' | cut -d ':' -f2 | xargs"),
		Arch:      executeCommand("hostnamectl | grep 'Architecture' | cut -d ':' -f2 | xargs"),
		DockerVer: executeCommand("docker -v | awk '{print $3}' | sed 's/,//'"),
		CPU:       executeCommand("top -bn1 | grep 'Cpu(s)' | sed 's/.*, *\\([0-9.]*\\)%* id.*/\\1/' | awk '{print 100 - $1 \"%\"}'"),
		Mem:       executeCommand("free | awk '/Mem:/ {printf \"%.2f%%\\n\", ($3-$6)/$2 * 100}'"),
	}
}
