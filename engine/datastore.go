package engine

import (
	"fmt"
	"log"
)

func (n *Node) VerifyDatastore() {
	// Check for broken/corrupt blocks in the datastore.
	broken, err := n.Node.Verify()
	if err != nil {
		log.Println(err)
		return
	}

	// Loop through found broken blocks.
	for pin := range broken {
		if n.Verbose {
			fmt.Printf("Corrupt Block Found: %s, re-pinning data from another node.\n", pin)
		}

		// Remove corrupt block from datastore.
		err = n.Node.Unpin(pin)
		if err != nil {
			log.Println(err)
			continue
		}

		// Re-pin data from another node.
		err = n.Node.Pin(pin)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
