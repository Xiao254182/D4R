package ps

import (
	"os/exec"
	"regexp"
	"strings"
)

// GetDockerContainers 获取当前运行中的 Docker 容器
func GetDockerContainers() ([][]string, error) {
	return getDockerOutput("docker", "ps")
}

// GetRunningContainerIDs 获取当前运行中的 Docker 容器的 ID
func GetRunningContainerIDs() ([]string, error) {
	output, err := exec.Command("docker", "ps", "-q").Output()
	if err != nil {
		return nil, err
	}
	return parseContainerIDs(output), nil
}

// getDockerOutput 通用函数，执行给定的 Docker 命令并解析结果
func getDockerOutput(command string, args ...string) ([][]string, error) {
	output, err := exec.Command(command, args...).Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var containers [][]string
	re := regexp.MustCompile(`\s{2,}`)

	for _, line := range lines[1:] { // 跳过表头
		if line = strings.TrimSpace(line); line != "" {
			fields := re.Split(line, -1)
			containers = append(containers, fields)
		}
	}

	return containers, nil
}

// parseContainerIDs 解析运行中的容器 ID
func parseContainerIDs(output []byte) []string {
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var containerIDs []string

	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			containerIDs = append(containerIDs, line)
		}
	}

	return containerIDs
}
