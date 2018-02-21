package util

import (
	"fmt"
	"os"
)

func HandleError(err error, fatal bool) {
	fmt.Println(err)
	if fatal {
		os.Exit(-1)
	}
}
