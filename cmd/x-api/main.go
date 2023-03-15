package main

import (
	"github.com/edalmi/x-api/commands"
)

func main() {
	if err := commands.New().Execute(); err != nil {
		panic(err)
	}
}
