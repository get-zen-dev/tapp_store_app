package k8sinterface

import (
	"errors"
	"net/url"
	"os/exec"
	"strings"
)

const commandCore = "microk8s"

type microk8sClient struct {
	domainName         string
	currentStatusCache string
}

func (m *microk8sClient) Start() error {
	if !checkInstallMicrok8s() {
		err := kuberInitialization()
		if err != nil {
			return err
		}
	} else if !checkSetupRepositoryOfAddons() {
		setupRepositoryOfAddons()
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
	return m.GetCachedModuleInfo(name)
}

func (m *microk8sClient) RemoveModule(name string) error {
	err := invokeCommand(commandCore, "disable", name)
	if err != nil {
		return err
	}
	return nil
}

func (m *microk8sClient) GetCachedModuleInfo(name string) (*ModuleInfo, error) {
	if m.currentStatusCache == "" {
		err := m.RefreshInfoCache()
		if err != nil {
			return nil, err
		}
	}

	status := m.currentStatusCache
	_, enableAndDisable, find := strings.Cut(status, "  enabled:")
	if !find {
		panic(status)
		return nil, errors.New(name + ": enable modules not found")
	}
	enable, disable, find := strings.Cut(enableAndDisable, "  disabled:")
	if !find {
		panic(status)
		return nil, errors.New(name + ": disabled modules not found")
	}

	isEnabled := strings.Index(enable, name)
	isDisabled := strings.Index(disable, name)

	if isEnabled == -1 && isDisabled == -1 {
		return nil, errors.New(name + ": module is not provided")
	}

	result := ModuleInfo{}
	result.Name = name
	result.IsEnabled = isEnabled != -1
	return &result, nil
}

func (m *microk8sClient) RefreshInfoCache() error {
	cmd := exec.Command(commandCore, "status")
	stdout, err := cmd.Output()
	m.currentStatusCache = string(stdout)
	return err
}

func kuberInitialization() error {
	err := invokeCommand("snap", "install", "microk8s", "--classic", "--channel=1.24/stable")
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

func (m *microk8sClient) GetModuleUrl(name string) url.URL {
	return url.URL{Host: m.domainName, Scheme: "https", Path: name}
}
