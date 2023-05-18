package k8sinterface

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const commandCore = "microk8s"

type microk8sClient struct {
}

func invokeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func (m *microk8sClient) Start() error {
	err := kuberInitialization()
	if err != nil {
		return err
	}

	command := "start"
	err = invokeCommand(commandCore, command)
	if err != nil {
		return err
	}

	return invokeCommand(commandCore, "status", "--wait-ready")
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
		if !checkInstallMicrok8s() {
			return errors.Join(errors.New("incorrect microk8s install"), err)
		}
	}

	addons := []string{"dns", "community", "traefik"}

	for i := range addons {
		addon := addons[i]
		err = invokeCommand(commandCore, "enable", addon)
		if err != nil {
			return errors.Join(errors.New("incorrect install of microk8s addon - "+addon), err)
		}
	}

	return nil
}

func checkInstall(app string) bool {
	cmd := fmt.Sprintf("dpkg -l | grep -w %s", app)
	_, err := exec.Command("sh", "-c", cmd).Output()
	return err == nil
}

func checkInstallMicrok8s() bool {
	cmd := "microk8s"
	_, err := exec.Command("sh", "-c", cmd).Output()
	return err != nil
}

func checkInstallSnap() bool {
	return checkInstall("snap")
}

func checkInGroupMicrok8s() bool {
	cmd := "groups | grep -w microk8s"
	_, err := exec.Command("sh", "-c", cmd).Output()
	return err == nil
}

func getUserName() (string, error) {
	cmd := "whoami"
	ans, err := exec.Command("sh", "-c", cmd).Output()
	return string(ans), err
}

func addGroupMicrok8s() error {
	username, err := getUserName()
	if err != nil {
		return err
	}
	cmd := fmt.Sprintf("sudo usermod -a -G microk8s %s", username)
	err = invokeCommand(cmd)
	if err != nil {
		return err
	}
	cmd = "newgrp microk8s"
	err = invokeCommand(cmd)
	if err != nil {
		return err
	}
	cmd = "sudo mkdir ~/.kube"
	err = invokeCommand(cmd)
	if err != nil {
		return err
	}
	cmd = fmt.Sprintf("sudo chown -R %s ~/.kube", username)
	err = invokeCommand(cmd)
	return err
}
