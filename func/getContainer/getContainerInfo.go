package getcontainer

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateContainerOut(containerName string) *tview.TextView {
	infoPanel := tview.NewTextView()
	infoPanel.SetBorder(true).SetTitle("Container Info").SetBorderColor(tcell.ColorLightSkyBlue)

	// 获取容器的详细信息
	info := CreateContainerInfo(containerName)
	infoPanel.SetText(info)

	return infoPanel
}

func CreateContainerInfo(containerName string) string {
	containerInfo := GetContainerInfo(containerName)
	mountsInfo := formatMounts(containerName)
	portsInfo := formatPorts(containerName)
	envInfo := formatEnvs(containerName)

	// 拼接要显示的信息
	info := fmt.Sprintf("状态: %s\n\n镜像: %s\n\n创建时间: %s\n\n挂载目录: \n%v\n端口映射: %s\n\n网络地址: %s\n\n工作目录: %s\n\n用户: %s\n\n环境变量: \n%s",
		string(containerInfo.State.Status), string(containerInfo.Config.Image), string(containerInfo.Created), mountsInfo, portsInfo, string(containerInfo.NetworkSettings.IPAddress), string(containerInfo.Config.WorkingDir), string(containerInfo.Config.User), envInfo)

	return info
}

func GetContainerInfo(containerName string) container.InspectResponse {
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
	return containerInfo
}

func formatMounts(containerName string) string {
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

	Mount := containerInfo.Mounts
	MountMap := make(chan string, len(Mount)) // 缓冲区大小为 Mount 切片的长度
	// 向通道中写入数据
	go func() {
		for _, mountMap := range Mount {
			MountMap <- fmt.Sprintf("%s:%s\n", mountMap.Source, mountMap.Destination)
		}
		close(MountMap) // 写入完成后关闭通道
	}()
	// 读取通道中的数据并拼接成字符串
	var mounts []string
	for mount := range MountMap {
		mounts = append(mounts, mount)
	}
	mountsInfo := ""
	for _, mount := range mounts {
		mountsInfo += mount
	}
	return mountsInfo
}

func formatPorts(containerName string) string {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.47"))
	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		return ""
	}
	defer cli.Close()

	containerID := containerName // 替换为具体的容器 ID 或名称
	containerInfo, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		fmt.Println("Error inspecting container:", err)
		return ""
	}

	// 获取容器的端口映射信息
	portMap := containerInfo.NetworkSettings.Ports

	var portsInfo string
	// 处理并格式化端口映射信息
	for port, bindings := range portMap {
		if len(bindings) > 0 {
			// 打印容器的端口
			portsInfo += fmt.Sprintf("%s，\n", port)

			// 用一个变量来存储绑定信息，最终输出时格式化为一行
			var bindInfo string
			for _, binding := range bindings {
				if binding.HostIP != "" && binding.HostPort != "" {
					bindInfo += fmt.Sprintf("%s:%s->%s, ", binding.HostIP, binding.HostPort, port)
				} else if binding.HostPort != "" {
					// 使用 [::] 替换 :: 以保持格式一致
					bindInfo += fmt.Sprintf("[::]:%s->%s, ", binding.HostPort, port)
				}
			}

			// 去掉最后一个多余的逗号和空格
			if len(bindInfo) > 0 && bindInfo[len(bindInfo)-2:] == ", " {
				bindInfo = bindInfo[:len(bindInfo)-2]
			}
			portsInfo += bindInfo + "\n"
		}
	}

	// 如果没有端口映射信息
	if portsInfo == "" {
		portsInfo = "没有端口映射信息"
	}

	// 去除最后的换行符，确保不多余的换行
	if len(portsInfo) > 0 && portsInfo[len(portsInfo)-1] == '\n' {
		portsInfo = portsInfo[:len(portsInfo)-1]
	}

	// 确保按指定格式返回，并且端口映射为多行
	return "\n" + portsInfo
}

func formatEnvs(containerName string) string {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.47"))
	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		return ""
	}
	defer cli.Close()

	containerID := containerName // 替换为具体的容器 ID 或名称
	containerInfo, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		fmt.Println("Error inspecting container:", err)
		return ""
	}

	// 获取容器的环境变量
	envVars := containerInfo.Config.Env

	// 格式化输出环境变量，去除方括号，并以逗号分隔
	return strings.Join(envVars, ", ")

}
