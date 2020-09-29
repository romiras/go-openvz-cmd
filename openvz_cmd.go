package openvzcmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type (
	// Options defines map of options
	Options map[string]string

	// CmdResult defines result of command
	CmdResult error

	ContainerInfo struct{}

	// IPOCCommander defines interface for management of OpenVZ
	IPOCCommander interface {
		ReadCommandsFromConfig(path string) error
		ListContainers() ([]ContainerInfo, CmdResult)
		CreateContainer(name, osTemplate string, options Options) CmdResult
		DeleteContainer(name string) CmdResult
		SetContainerParameters(options Options) CmdResult
		StartContainer(name string, wait bool) CmdResult
		StopContainer(name string, wait bool) CmdResult
		RestartContainer(name string) CmdResult
		SuspendContainer(name string) CmdResult
		ResumeContainer(name string) CmdResult
	}
)

type ExecCommandInfo struct {
	Program   string   `yaml:"program"`
	Arguments []string `yaml:"arguments"`
	Vars      []string `yaml:"vars"`
}

type ExecCommandsMap struct {
	CtCreate ExecCommandInfo `yaml:"ct-create"`
	CtSet    ExecCommandInfo `yaml:"ct-set"`
	CtDelete ExecCommandInfo `yaml:"ct-delete"`
}

type POCCommanderStub struct {
	execCommandsMap ExecCommandsMap
}

func (cmd *POCCommanderStub) readCommandsFromConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, &cmd.execCommandsMap)
	if err != nil {
		return err
	}

	return nil
}

func buildCommand(execInfo *ExecCommandInfo, params Options) *exec.Cmd {
	fmt.Printf("%+v\n", execInfo)

	args := bindVars(execInfo, params)
	fmt.Println(args)

	command := exec.Command(execInfo.Program, args...)
	return command
}

func bindVars(execInfo *ExecCommandInfo, params Options) []string {
	args := make([]string, 0)

	for _, arg := range execInfo.Arguments {
		a := arg
		for _, v := range execInfo.Vars {
			if value, ok := params[v]; ok {
				a = strings.ReplaceAll(a, "{{"+v+"}}", value)
			}
		}
		args = append(args, a)
	}

	return args
}

func (cmd *POCCommanderStub) CreateContainer(name, osTemplate string, options Options) CmdResult {
	params := make(Options)
	params["name"] = name
	params["ostemplate"] = osTemplate
	command := buildCommand(&cmd.execCommandsMap.CtCreate, params)

	return command.Run()
}

func (cmd *POCCommanderStub) SetContainerParameters(name string, params Options) CmdResult {
	if params == nil {
		log.Fatal("params cannot be nil")
	}
	params["name"] = name
	return cmd.execCommand(&cmd.execCommandsMap.CtSet, params)
}

func (cmd *POCCommanderStub) DeleteContainer(name string) CmdResult {
	return cmd.execCommand(&cmd.execCommandsMap.CtDelete, nil)
}

func (cmd *POCCommanderStub) execCommand(execInfo *ExecCommandInfo, params Options) CmdResult {
	command := buildCommand(execInfo, params)
	var out bytes.Buffer
	command.Stdout = &out
	fmt.Printf("Stdout:\n%s\n", out.String())
	return command.Run()
}

func NewPOCCommanderStub(path string) (*POCCommanderStub, error) {
	cmd := &POCCommanderStub{}
	err := cmd.readCommandsFromConfig(path)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}
