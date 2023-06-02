package environment

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	owner                        = "get-zen-dev"
	repository                   = "tapp_store_rep"
	path                         = "addons"
	ref                          = "main"
	domain                       = ""
	currentVersions *viper.Viper = nil
)

// Returns the name of owner of the repository with addons
func GetOwner() string {
	return owner
}

// Returns the name of the repository with addons
func GetRepository() string {
	return repository
}

// Returns the path to the addons in the repository
func GetPath() string {
	return path
}

// Returns the ref in the repository
func GetRef() string {
	return ref
}

// Returns the domain
func GetDomain() (string, error) {
	if domain != "" {
		return domain, nil
	}
	domainRead, err := ReadFromConfig("app.yaml", "domain")
	domain = domainRead
	return domainRead, err
}

func ReadFromConfigCurrentVersion(key string) (string, error) {
	if currentVersions == nil {
		currentVersions = initViper("current_version.yaml")
	}
	data := currentVersions.Get(key)
	switch data.(type) {
	case string:
		return data.(string), nil
	default:
		return "", fmt.Errorf("not found")
	}
}

func WriteInConfigCurrentVersion(key, value string) error {
	if currentVersions == nil {
		currentVersions = initViper("current_version.yaml")
	}
	currentVersions.Set(key, value)
	err := currentVersions.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func initViper(file string) *viper.Viper {
	v := viper.New()
	v.SetConfigFile("./../configs/" + file)
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			panic(err)
		}
	}
	return v
}

func WriteInConfig(file, key, value string) error {
	v := initViper(file)
	v.Set(key, value)
	err := v.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func ReadFromConfig(file, key string) (string, error) {
	v := initViper(file)
	data := v.Get(key)
	switch data.(type) {
	case string:
		return data.(string), nil
	default:
		return "", fmt.Errorf("not found")
	}
}

func ReadInfoAddonsModels() *Models {
	slice := ReadInfoAddonsSlice()
	models := Models{}
	for _, v := range *slice {
		models.append(Model{v["name"], v["version"], v["description"]})
	}
	return &models
}

func ReadInfoAddonsSlice() *[]map[string]string {
	viper.SetConfigFile("./../configs/addons.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	slice := []map[string]string{}
	viper.UnmarshalKey("microk8s-addons.addons", &slice)
	return &slice
}
