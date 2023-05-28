package k8sinterface

import (
	"errors"
	"net/url"
)

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
	GetCachedModuleInfo(name string) (*ModuleInfo, error)
	RefreshInfoCache() error
	GetModuleUrl(name string) (url.URL, error)
}

func GetInterfaceProvider(domain string) (KuberInterface, error) {
	if !checkIsRootGranted() {
		return nil, errors.New("cannot user interface without root privileges")
	}

	return &microk8sClient{domain, ""}, nil
}
