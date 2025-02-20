package getcontainer

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateContainerOut(containerName string) *tview.TextView {
	infoPanel := tview.NewTextView()
	infoPanel.SetBorder(true).SetTitle("Container Info").SetBorderColor(tcell.ColorLightSkyBlue)

	// 获取容器的详细信息
	info := GetContainerInfo(containerName)
	infoPanel.SetText(info)

	return infoPanel
}

func GetContainerInfo(containerName string) string {
	cmd := exec.Command("docker", "inspect", containerName)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("无法获取容器信息: %v", err)
	}

	// 获取容器状态
	statusCmd := exec.Command("docker", "inspect", "--format", "{{.State.Status}}", containerName)
	status, err := statusCmd.Output()
	if err != nil {
		return fmt.Sprintf("容器状态获取失败: %v", err)
	}

	// 获取容器镜像
	imageCmd := exec.Command("docker", "inspect", "--format", "{{.Config.Image}}", containerName)
	image, err := imageCmd.Output()
	if err != nil {
		return fmt.Sprintf("容器镜像获取失败: %v", err)
	}

	// 获取容器创建时间
	createdCmd := exec.Command("docker", "inspect", "--format", "{{.Created}}", containerName)
	created, err := createdCmd.Output()
	if err != nil {
		return fmt.Sprintf("容器创建时间获取失败: %v", err)
	}

	// 转换容器创建时间为中国时间
	dateCmd := exec.Command("bash", "-c", fmt.Sprintf(`export TZ="Asia/Shanghai"; date -d "%s" "+%%Y-%%m-%%d %%H:%%M:%%S"`, strings.TrimSpace(string(created))))
	createdate, err := dateCmd.Output()
	if err != nil {
		return fmt.Sprintf("容器创建时间转换失败: %v", err)
	}

	// 获取容器的挂载目录 (-v)
	volumesCmd := exec.Command("docker", "inspect", "--format", "{{range .Mounts}}{{.Source}}:{{.Destination}}\n{{end}}", containerName)
	volumes, err := volumesCmd.Output()
	if err != nil {
		return fmt.Sprintf("容器挂载目录获取失败: %v", err)
	}

	// 获取容器的端口映射 (-p)
	portsCmd := exec.Command("docker", "ps", "-a", "--filter", "name=^"+containerName+"$", "--format", "table {{.Ports}}")
	ports, err := portsCmd.Output()
	if err != nil {
		return fmt.Sprintf("容器端口映射获取失败: %v", err)
	}

	// 获取容器的网络配置 (--network)
	networkCmd := exec.Command("docker", "inspect", "--format", "{{.NetworkSettings.IPAddress}}", containerName)
	network, err := networkCmd.Output()
	if err != nil {
		return fmt.Sprintf("容器网络地址获取失败: %v", err)
	}

	// 获取容器的工作目录 (-w)
	workingDirCmd := exec.Command("docker", "inspect", "--format", "{{.Config.WorkingDir}}", containerName)
	workingDir, err := workingDirCmd.Output()
	if err != nil {
		return fmt.Sprintf("容器工作目录获取失败: %v", err)
	}

	// 获取容器的用户 (-u)
	userCmd := exec.Command("docker", "inspect", "--format", "{{.Config.User}}", containerName)
	user, err := userCmd.Output()
	if err != nil {
		return fmt.Sprintf("容器用户获取失败: %v", err)
	}

	// 获取容器的环境变量 (-e)
	envCmd := exec.Command("docker", "inspect", "--format", "{{range .Config.Env}}{{.}}{{end}}", containerName)
	env, err := envCmd.Output()
	if err != nil {
		return fmt.Sprintf("容器环境变量获取失败: %v", err)
	}

	// 拼接要显示的信息
	info := fmt.Sprintf("状态: %s\n镜像: %s\n创建时间: %s\n挂载目录: \n%s\n端口映射: %s\n网络地址: %s\n工作目录: %s\n用户: %s\n环境变量: \n%s",
		string(status), string(image), string(createdate), string(volumes), string(ports), string(network), string(workingDir), string(user), string(env))

	return info
}
