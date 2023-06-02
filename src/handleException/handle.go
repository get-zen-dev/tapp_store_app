package handleexception

import (
	"fmt"
	"os"
)

func PrintErr(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func PrintErrorIfNotNil(err error) {
	if err != nil {
		PrintErr(err)
	}
}
