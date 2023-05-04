package k8sinterface

type ModuleInfo struct {
	Name    string
	Status  string
	Route   string
	Version string
}

type KuberInterface interface {
	InitKuber()

	InstallModule(name string) ModuleInfo
	RemoveModule(name string) bool
	GetModuleInfo(name string) ModuleInfo
	//UpdateModule(name string) ModuleInfo
}
