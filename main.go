package main

import (
	"os"

	"github.com/youngminz/appstore-scraper-cli/cmd"
)

func main() {
	if err := cmd.Execute(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		os.Exit(1)
	}
}
