package main

import (
	"fmt"
	"log"
	"time"

	"github.com/arkenproject/arken/ipfs"

	"github.com/arkenproject/arken/engine"
	"github.com/arkenproject/arken/stats"

	"github.com/alecthomas/units"
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

		// Generate Stats Data
		n, err := units.ParseBase2Bytes(config.Global.General.PoolSize)
		if err != nil {
			log.Fatal(err)
		}

		total := int(n)

		usage, err := ipfs.GetRepoSize()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(usage)
		data := stats.NodeData{
			ID:         ipfs.GetID(),
			Username:   config.Global.Stats.Username,
			Email:      config.Global.Stats.Email,
			TotalSpace: total / 1000000000,
			UsedSpace:  usage / 1000000000,
		}

		stats.CheckIn("https://arken.io/beacon", data)

		fmt.Println("\n[System Sleeping for 1 Hour]")

		time.Sleep(1 * time.Hour)
	}
}
