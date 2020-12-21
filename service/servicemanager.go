package service

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

type Supervisor int

const (
	Perp  Supervisor = iota
	Runit Supervisor = iota
)

type Service struct {
	GetServiceStatus (ServiceStatus)
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func Restart(supervisor Supervisor, serviceName string) {
	switch supervisor {
	case Perp:
		execCommandWithEnv("perpctl", "d", serviceName)
		execCommandWithEnv("perpctl", "u", serviceName)
	case Runit:
		execCommandWithEnv("/mnt/sdcard/hacks/runit/bin/sv", "restart", serviceName)
	}
}

func Stop(supervisor Supervisor, serviceName string) {
	switch supervisor {
	case Perp:
		execCommandWithEnv("perpctl", "d", serviceName)
	case Runit:
		execCommandWithEnv("/mnt/sdcard/hacks/runit/bin/sv", "down", serviceName)
	}
}

func Status(friendlyName string, supervisor Supervisor, serviceName string) ServiceStatus {
	status := ServiceStatus{Name: friendlyName, Status: "unknown"}

	switch supervisor {
	case Runit:
		output := execCommandWithEnv("sv", "status", serviceName)
		if strings.Contains(output, ":") {
			status = ServiceStatus{
				Name:   friendlyName,
				Status: output[0:strings.Index(output, ":")],
			}
		}
	}

	return status
}

func execCommandWithEnv(name string, arg ...string) string {
	var outputBuffer, errorBuffer bytes.Buffer

	command := exec.Command(name, arg...)
	command.Env = append(os.Environ(), "SVDIR=/mnt/data/etc/runit")
	command.Stdout = &outputBuffer
	command.Stderr = &errorBuffer
	command.Run()

	return outputBuffer.String()
}
