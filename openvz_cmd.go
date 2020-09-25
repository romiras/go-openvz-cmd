package openvzcmd

import (
	"os/exec"
)

type (
	// Options defines map of options
	Options map[string]string

	// CmdResult defines result of command
	CmdResult exec.ExitError

	ContainerInfo struct{}

	// IPOCCommander defines interface for management of OpenVZ
	IPOCCommander interface {
		ReadCommandsFromConfig(path string) error
		ListContainers() ([]ContainerInfo, CmdResult)
		CreateContainer(name, osTemplate string, options ...Options) CmdResult
		DeleteContainer(name string) CmdResult
		SetContainerParameters(options Options) CmdResult
		StartContainer(name string, wait bool) CmdResult
		StopContainer(name string, wait bool) CmdResult
		RestartContainer(name string) CmdResult
		SuspendContainer(name string) CmdResult
		ResumeContainer(name string) CmdResult
	}
)
