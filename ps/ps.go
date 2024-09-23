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

			// 确保 fields 至少有 6 个元素
			if len(fields) < 7 {
				// 如果字段数不足 6，补充空值
				for len(fields) < 7 {
					fields = append(fields, "")
				}
				// 检查第六个值并相应调整第五和第六个值
				if fields[6] == "" {
					// 第六个值为空，将第六个值赋给第五个值，保持第五个值为空
					fields[6] = fields[5]
					fields[5] = ""
				}
			}

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

// GetDockerCompose 获取当前运行中的 Docker Compose 集群
func GetDockerCompose() ([][]string, error) {
	return getDockerComposeOutput("docker-compose", "ls")
}

// GetRunningContainerIDs 获取当前运行中的 Docker 容器的 ID
func GetRunningDockerComposeName() ([]string, error) {
	output, err := exec.Command("docker-compose", "ls", "-q").Output()
	if err != nil {
		return nil, err
	}
	return parseContainerIDs(output), nil
}

// getDockerComposeOutput 通用函数，执行给定的 Docker 命令并解析结果
func getDockerComposeOutput(command string, args ...string) ([][]string, error) {
	output, err := exec.Command(command, args...).Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var dockerCompose [][]string
	re := regexp.MustCompile(`\s{2,}`)

	for _, line := range lines[1:] { // 跳过表头
		if line = strings.TrimSpace(line); line != "" {
			fields := re.Split(line, -1)
			dockerCompose = append(dockerCompose, fields)
		}
	}

	return dockerCompose, nil
}
