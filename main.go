package main

import (
	"os"

	"github.com/redhat-developer/opencompose/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
