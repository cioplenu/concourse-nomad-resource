package common

import (
	"fmt"
	"os"
)

func Check(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg+": %s\n", err)
		os.Exit(1)
	}
}
