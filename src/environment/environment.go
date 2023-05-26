package environment

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	owner      = "canonical"
	repository = "microk8s-core-addons"
	path       = "addons"
	ref        = "main"
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

func WriteInConfig(key, value string) error {
	viper.SetConfigFile("./../configs/app.env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	viper.Set(key, value)
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func ReadFromConfig(key string) (string, error) {
	viper.SetConfigFile("./../configs/app.env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	data := viper.Get(key)
	switch data.(type) {
	case string:
		return data.(string), nil
	default:
		return "", fmt.Errorf("domen not found")
	}
}

func ReadInfoAddonsMap() *map[string]map[string]string {
	slice := ReadInfoAddonsSlice()
	info := map[string]map[string]string{}
	for _, v := range *slice {
		info[v["name"]] = map[string]string{}
		info[v["name"]]["version"] = v["version"]
		info[v["name"]]["description"] = v["description"]
	}
	return &info
}

func ReadInfoAddonsSlice() *[]map[string]string {
	viper.SetConfigFile("./../configs/addons.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	slice := []map[string]string{}
	viper.UnmarshalKey("microk8s-addons.addons", &slice)
	return &slice
}
