package main

import (
	"fmt"
	"log"
	"time"

	"github.com/arkenproject/arken/engine"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/keysets"
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

		err := keysets.LoadSets(config.Keysets)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\n[Starting Rebalancing]")

		err = engine.Rebalance()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\n[Finished Data Rebalance]")

		fmt.Println("\n[System Sleeping for 1 Hour]")

		time.Sleep(1 * time.Hour)
	}
}
