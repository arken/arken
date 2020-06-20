package main

import (
	"fmt"

	"github.com/archivalists/arken/engine"

	"github.com/archivalists/arken/config"
	"github.com/archivalists/arken/keysets"
)

func main() {
	fmt.Println("Welcome to Arken!")
	fmt.Printf("Application Version %s\n\n", config.Global.General.Version)

	fmt.Println("Arken is now in [System Startup]")
	keysets.LoadSets(config.Keysets.Sets)

	engine.Rebalance()
}
