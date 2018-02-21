package util

import (
	"fmt"
	"os"
)

// HandleError is a function for handling errors
func HandleError(err error, fatal bool) {
	fmt.Println(err)
	if fatal {
		os.Exit(-1)
	}
}
