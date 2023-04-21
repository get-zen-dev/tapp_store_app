package environment

import (
	"fmt"

	"github.com/spf13/viper"
)

var owner string
var repository string
var path string
var ref string

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
