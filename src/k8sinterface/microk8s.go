package k8sinterface

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"strings"
	"time"
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
		err := setupRepositoryOfAddons()
		if err != nil {
			return err
		}
	}
	err := invokeCommand(commandCore, "start")
	if err != nil {
		return err
	}
	for i := 0; i < 5 && !strings.Contains(m.currentStatusCache, "microk8s is running"); i++ {
		m.RefreshInfoCache()
		time.Sleep(time.Second)
	}
	if !strings.Contains(m.currentStatusCache, "microk8s is running") {
		return fmt.Errorf("failed to start microk8s client")
	}
	return nil
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
		return nil, errors.New("enable modules not found")
	}
	enable, disable, find := strings.Cut(enableAndDisable, "  disabled:")
	if !find {
		return nil, errors.New("disabled modules not found")
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
	out, err := invokeCommandWithOutput("snap", "install", "microk8s", "--classic", "--channel=1.25/stable")
	if err != nil {
		return errors.Join(errors.New("incorrect microk8s install: "+out), err)
	}

	out, err = invokeCommandWithOutput(commandCore, "start")
	if err != nil {
		return errors.Join(errors.New("error on microk8s starting: "+out), err)
	}

	toInvoke := []string{"enable"}
	toInvoke = append(toInvoke, microk8s_addons...)

	out, err = invokeCommandWithOutput(commandCore, toInvoke...)
	if err != nil {
		return errors.Join(errors.New("error on microk8s addons installing: "+out), err)
	}

	for i := 0; i < 5 && !strings.Contains(out, "microk8s is running"); i++ {
		err = invokeCommand(commandCore, "status", "--wait-ready")
		println(i)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to start microk8s client")
	}

	println("SUPER MEGA WAITING:")
	time.Sleep(time.Second * 30)
	println("Installing cert:")
	err = applyConfigToKuber(cfg_certs_yml)
	if err != nil {
		return errors.Join(errors.New("error on microk8s installing certs"), err)
	}

	return setupRepositoryOfAddons()
}

func (m *microk8sClient) GetModuleUrl(name string) url.URL {
	return url.URL{Host: name + "." + m.domainName, Scheme: "https"}
}

func applyConfigToKuber(config string) error {

	cmd := exec.Command(commandCore, "kubectl", "apply", "-f", "-")
	reader, writer := io.Pipe()
	cmd.Stdin = reader

	go func() {
		defer writer.Close()
		writer.Write([]byte(config))
	}()

	out, err := cmd.Output()
	if err != nil {
		return errors.Join(errors.New("Incorrect apply certs: "+string(out)), err)
	}

	return nil
}
