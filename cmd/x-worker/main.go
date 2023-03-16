package main

import (
	"github.com/edalmi/x-api/internal/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		panic(err)
	}
}
