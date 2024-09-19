package rm

import (
	"os/exec"
)

// DeleteContainer 删除指定 ID 的容器
func DeleteContainer(containerID string) error {
	// 执行 Docker 删除命令
	cmd := exec.Command("docker", "rm", "-f", containerID)
	return cmd.Run()
}
