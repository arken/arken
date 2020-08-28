package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/arkenproject/arken/engine"
	"github.com/arkenproject/arken/stats"

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

	// Check whether to report node stats
	if strings.ToLower(config.Global.General.StatsReporting) == "on" {
		go func() {
			for {
				// If allowed report the stats to the keyset stats server.
				err := stats.Report(config.Keysets)
				if err != nil {
					log.Println(err)
				}
				time.Sleep(1 * time.Hour)
			}
		}()
	}

	// Verify Locally Pinned Files and Re-Pin if lost.
	go func() {
		for {
			engine.VerifyLocal()
			time.Sleep(1 * time.Hour)
		}
	}()

	for {
		fmt.Println("\n[Indexing & Updating Keysets]")

		err := keysets.LoadSets(config.Keysets)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\n[Starting Rebalancing]")
		hit, err := engine.CheckNetUsage()
		if err != nil {
			log.Fatal(err)
		}

		if hit {
			fmt.Println("[Cancelling Rebalance due to Network Limit Hit]")
		} else {
			err = engine.Rebalance()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("\n[Finished Data Rebalance]")
		}
		fmt.Println("\n[System Sleeping for 1 Hour]")

		time.Sleep(1 * time.Hour)
	}
}
