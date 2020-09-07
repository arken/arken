package main

import (
	"fmt"
	"strings"

	"github.com/arkenproject/arken/tasks"

	"github.com/arkenproject/arken/ipfs"

	"github.com/arkenproject/arken/config"
)

func main() {
	fmt.Println(`                                            
    // | |                                       
   //__| |     __     / ___      ___       __    
  / ___  |   //  ) ) //\ \     //___) ) //   ) ) 
 //    | |  //      //  \ \   //       //   / /  
//     | | //      //    \ \ ((____   //   / /   `)

	fmt.Printf("Application Version %s\n\n", config.Global.General.Version)

	fmt.Println("Arken is now in [System Startup]")
	ipfs.Init()

	// Launch Stats Reporting if enabled in the config.
	if strings.ToLower(config.Global.General.StatsReporting) == "on" {
		go tasks.StatsReporting()
	}

	// Verify Locally Pinned Files and Re-Pin if lost.
	go tasks.VerifyLocal()

	// Begin the main Arken process.
	tasks.Main()
}
