package k8sinterface

import "errors"

type ModuleInfo struct {
	Name      string
	IsEnabled bool
	//Route   string
	//Version string
}

type KuberInterface interface {
	Start() error
	Stop() error

	InstallModule(name string) (*ModuleInfo, error)
	RemoveModule(name string) error
	GetModuleInfo(name string) (*ModuleInfo, error)
	//UpdateModule(name string) ModuleInfo
}

func GetInterfaceProvider(domain string) (KuberInterface, error) {
	if !checkIsRootGranted() {
		return nil, errors.New("cannot user interface without root privileges")
	}

	return &microk8sClient{domain}, nil
}
