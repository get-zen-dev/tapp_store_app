package k8sinterface

import (
	"errors"
	"fmt"
	"os"
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

func invokeCommandInteractive(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func InitKuberInteractive() error {
	installMicrok8s := "snap install microk8s --classic"
	cmd := exec.Command("sh", "-c", installMicrok8s)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	addons := []string{"dns", "community", "traefik"}

	for i := range addons {
		installAddon := fmt.Sprintf("%s enable %s", commandCore, addons[i])
		cmd := exec.Command("sh", "-c", installAddon)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func CheckInstallMicrok8s() bool {
	cmd := fmt.Sprintf("dpkg -l | grep -w %s", commandCore)
	_, err := exec.Command("sh", "-c", cmd).Output()
	return err == nil
}
func CheckInstallSnap() bool {
	cmd := "dpkg -l | grep -w snap"
	_, err := exec.Command("sh", "-c", cmd).Output()
	return err == nil
}

func InstallSnap() error {
	update := "sudo apt update"
	installSnap := "sudo apt install snapd"
	cmd := exec.Command("sh", "-c", update)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(err.Error(), errors.New("error update"))
	}
	cmd = exec.Command("sh", "-c", installSnap)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf(err.Error(), errors.New("error install"))
	}
	return nil
}

var su bool = checkGranted()

func checkGranted() bool {
	cmd := "id"
	out, _ := exec.Command("sh", "-c", cmd).Output()
	return strings.Contains(string(out), "uid=0")
}

func SUStatus() bool {
	return su
}

func CheckInstalledOrInstalKuber() error {
	var err error
	if !CheckInstallSnap() {
		fmt.Println("Snap package manager is not installed. You need to install. Enter password to install:")
		err = InstallSnap()
	}
	if err != nil {
		return err
	}
	if !CheckInstallMicrok8s() {
		fmt.Println("Microk8s is not installed.\nThe following addons will also be installed: dns, community, traefik. \nEnter password to install:")
		err = InitKuberInteractive()
	}
	if err != nil {
		return err
	}
	return nil
}
