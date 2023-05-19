package environment

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	owner      string
	repository string
	path       string
	ref        string
)

func SetUpEnv() {
	viper.SetConfigFile("./../configs/app.env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	owner = viper.Get("owner").(string)
	repository = viper.Get("repository").(string)
	path = viper.Get("path").(string)
	ref = viper.Get("ref").(string)
}

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
