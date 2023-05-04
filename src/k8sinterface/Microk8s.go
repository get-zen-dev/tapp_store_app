package k8sinterface

import (
	"errors"
	"os/exec"
	"strings"
)

const commandCore = "microk8s"

type Microk8sClient struct {
}

func invokeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func (m *Microk8sClient) InitKuber() error {
	err := invokeCommand("snap", "install", "microk8s", "--classic")
	if err != nil {
		return err
	}

	addons := []string{"dns", "community", "traefik"}

	for i := range addons {
		err = invokeCommand(commandCore, "enable", addons[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Microk8sClient) InstallModule(name string) (*ModuleInfo, error) {
	err := invokeCommand(commandCore, "enable", name)
	if err != nil {
		return nil, err
	}
	return m.GetModuleInfo(name)
}

func (m *Microk8sClient) RemoveModule(name string) error {
	err := invokeCommand(commandCore, "disable", name)
	if err != nil {
		return err
	}
	return nil
}

func (m *Microk8sClient) GetModuleInfo(name string) (*ModuleInfo, error) {
	cmd := exec.Command(commandCore, "enable", name)
	stdout, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	status := string(stdout)
	_, enableAndDisable, find := strings.Cut(status, "  enabled:")
	if find {
		return nil, errors.New("enable modules not found")
	}
	enable, disable, find := strings.Cut(enableAndDisable, "  disabled:")
	if find {
		return nil, errors.New("disabled modules not found")
	}

	isEnabled := strings.Index(enable, name)
	isDisabled := strings.Index(disable, name)

	if isEnabled == -1 && isDisabled == -1 {
		return nil, errors.New("module is not provided in kuber")
	}

	result := ModuleInfo{}
	result.Name = name
	result.IsEnabled = isEnabled != -1
	return &result, nil
}
