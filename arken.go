package main

import (
	"fmt"

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

	// Begin the main Arken process.
	tasks.Main()
}
