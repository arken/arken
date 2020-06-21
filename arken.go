package main

import (
	"fmt"
	"time"

	"github.com/archivalists/arken/engine"

	"github.com/archivalists/arken/config"
	"github.com/archivalists/arken/keysets"
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
	for {
		fmt.Println("\n[Indexing & Updating Keysets]")

		keysets.LoadSets(config.Keysets.Sets)

		fmt.Println("\n[Starting Rebalancing]")

		engine.Rebalance()

		fmt.Println("\n[Finished Data Rebalance]")

		fmt.Printf("\n[System Sleeping for 10 Seconds]")

		time.Sleep(10 * time.Second)
	}
}
