package k8sinterface

import (
	"errors"
	"os/exec"
	"strings"
)

const commandCore = "microk8s"

type microk8sClient struct {
	domainName string
}

func (m *microk8sClient) Start() error {
	if !checkInstallMicrok8s() {
		err := kuberInitialization()
		if err != nil {
			return err
		}
	}

	return invokeCommand(commandCore, "start")
}

func (m *microk8sClient) Stop() error {
	command := "stop"
	return invokeCommand(commandCore, command)
}

func (m *microk8sClient) InstallModule(name string) (*ModuleInfo, error) {
	err := invokeCommand(commandCore, "enable", name)
	if err != nil {
		return nil, err
	}
	return m.GetModuleInfo(name)
}

func (m *microk8sClient) RemoveModule(name string) error {
	err := invokeCommand(commandCore, "disable", name)
	if err != nil {
		return err
	}
	return nil
}

func (m *microk8sClient) GetModuleInfo(name string) (*ModuleInfo, error) {
	cmd := exec.Command(commandCore, "status")
	stdout, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	status := string(stdout)
	_, enableAndDisable, find := strings.Cut(status, "  enabled:")
	if !find {
		return nil, errors.New("enable modules not found")
	}
	enable, disable, find := strings.Cut(enableAndDisable, "  disabled:")
	if !find {
		return nil, errors.New("disabled modules not found")
	}

	isEnabled := strings.Index(enable, name)
	isDisabled := strings.Index(disable, name)

	if isEnabled == -1 && isDisabled == -1 {
		return nil, errors.New("module is not provided")
	}

	result := ModuleInfo{}
	result.Name = name
	result.IsEnabled = isEnabled != -1
	return &result, nil
}

func kuberInitialization() error {
	err := invokeCommand("snap", "install", "microk8s", "--classic")
	if err != nil {
		return errors.Join(errors.New("incorrect microk8s install"), err)
	}

	err = invokeCommand(commandCore, "start")
	if err != nil {
		return errors.Join(errors.New("error on microk8s starting"), err)
	}

	addons := []string{"dns", "community", "traefik"}

	for i := range addons {
		addon := addons[i]
		err = invokeCommand(commandCore, "enable", addon)
		if err != nil {
			return errors.Join(errors.New("incorrect install of microk8s addon - "+addon), err)
		}
	}

	return setupRepositoryOfAddons()
}
