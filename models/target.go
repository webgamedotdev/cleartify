package models

import (
	"fmt"
	"os/exec"
	"strings"
)

var Targets map[string]Target

func RegisterTarget(name string, target Target) {
	if Targets == nil {
		Targets = make(map[string]Target)
	}
	Targets[name] = target
}

func init() {
	RegisterTarget("local", &LocalTarget{})
	RegisterTarget("docker", &DockerTarget{})
}

type Target interface {
	ExecuteCommand(command string) (string, error)
}

// LocalTarget implements Target for local machine
type LocalTarget struct{}

func (lt *LocalTarget) ExecuteCommand(command string) (string, error) {
	// Split the command into command and arguments
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)

	// Run the command
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// DockerTarget implements Target for Docker containers
type DockerTarget struct {
	ContainerID string
}

func (dt *DockerTarget) ExecuteCommand(command string) (string, error) {
	// Construct the full docker exec command
	cmd := exec.Command("docker", "exec", dt.ContainerID, "bash", "-c", command)

	// Run the command
	output, err := cmd.CombinedOutput()
	fmt.Println("Output in DockerTarget ExecuteCommand:", string(output))
	return string(output), err
}
