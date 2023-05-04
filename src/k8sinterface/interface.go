package k8sinterface

type ModuleInfo struct {
	Name      string
	IsEnabled bool
	//Route   string
	//Version string
}

type KuberInterface interface {
	InitKuber() error

	InstallModule(name string) (*ModuleInfo, error)
	RemoveModule(name string) error
	GetModuleInfo(name string) (*ModuleInfo, error)
	//UpdateModule(name string) ModuleInfo
}
