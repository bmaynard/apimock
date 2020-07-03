package main

import (
	"fmt"
	"os"

	"github.com/bmaynard/apimock/cmd"
)

func main() {
	if err := cmd.NewApiMockApp().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
