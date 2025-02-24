package getcontainer

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
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
	cli, err := client.NewClientWithOpts(client.WithVersion("1.47"))
	if err != nil {
		fmt.Println("Error creating Docker client:", err)
	}
	defer cli.Close()

	containerID := containerName // 替换为具体的容器 ID 或名称
	containerInfo, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		fmt.Println("Error inspecting container:", err)
	}

	// 拼接要显示的信息
	info := fmt.Sprintf("状态: %s\n镜像: %s\n创建时间: %s\n挂载目录: \n%s\n端口映射: %s\n网络地址: %s\n工作目录: %s\n用户: %s\n环境变量: \n%s",
		string(containerInfo.State.Status), string(containerInfo.Config.Image), string(containerInfo.Created), containerInfo.Mounts, containerInfo.NetworkSettings.Ports, string(containerInfo.NetworkSettings.IPAddress), string(containerInfo.Config.WorkingDir), string(containerInfo.Config.User), containerInfo.Config.Env)

	return info
}
