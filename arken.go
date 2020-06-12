package main

import (
	"fmt"

	"github.com/archivalists/arken/config"
)

func main() {
	fmt.Println(config.Global.General.Version)
}
