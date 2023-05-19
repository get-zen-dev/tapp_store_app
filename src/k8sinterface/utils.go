package k8sinterface

import (
	"fmt"
	"os/exec"
	"strings"
)

func checkIsRootGranted() bool {
	cmd := "id"
	out, _ := exec.Command("sh", "-c", cmd).Output()
	return strings.Contains(string(out), "uid=0")
}

func invokeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	_, err := cmd.Output()
	return err
}

func checkInstall(app string) bool {
	cmd := fmt.Sprintf("dpkg -l | grep -w %s", app)
	_, err := exec.Command("sh", "-c", cmd).Output()
	return err == nil
}

func checkInstallMicrok8s() bool {
	err := invokeCommand("microk8s", "version")
	return err == nil
}

func getUserName() (string, error) {
	cmd := "whoami"
	ans, err := exec.Command("sh", "-c", cmd).Output()
	return string(ans), err
}

func checkInGroupMicrok8s() bool {
	cmd := "groups | grep -w microk8s"
	_, err := exec.Command("sh", "-c", cmd).Output()
	return err == nil
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
