package main

import (
	"fmt"
	"os"
)

const (
	appName = "Fingerguns"
	command = "fgg"
	version = "0.2.7"
)

var (
	validDirectivesFilenames = map[string]bool{
		"fgg.toml":        true,
		"fingerguns.toml": true,
		"scripts.toml":    true,
		"commands.toml":   true,
		"exes.toml":       true,
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
