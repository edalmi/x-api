package main

import (
	"github.com/edalmi/x-api/internal/cmd"
)

func main() {
	if err := cmd.NewCmdRoot().Execute(); err != nil {
		panic(err)
	}
}
