package engine

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"

	"github.com/arken/arken/config"
	"github.com/arken/arken/database"
	"github.com/arken/arken/ipfs"
	"github.com/arken/arken/manifest"
)

type Node struct {
	Cfg      *config.Config
	DB       *database.DB
	Node     *ipfs.Node
	Manifest *manifest.Manifest
	Verbose  bool
}

func (n *Node) FileAdder() (chan<- database.File, error) {
	input := make(chan database.File, 10)

	// Determine the possible number of threads for the system's CPU
	workers := genNumWorkers()

	// Generate Worker Threads
	for i := 0; i < workers; i++ {
		go n.addWorker(input)
	}

	return input, nil
}

func (n *Node) addWorker(input <-chan database.File) {
	for file := range input {
		// Check the number of times a file is replicated across the cluster.
		replications, err := n.Node.FindProvs(file.ID, int(n.Manifest.Replications))
		if err != nil {
			log.Println(err)
			continue
		}

		// If in verbose mode display print out of replications
		if n.Verbose {
			fmt.Printf("File: %s is backed up %d time(s) and the threshold is %d.\n",
				file.ID, replications, n.Manifest.Replications,
			)
		}

		// If file is replicated at least once but not enough times attempt to pin in
		// locally if activation energy is high enough.
		if replications < int(n.Manifest.Replications) && replications >= 1 {
			activationEnergy := float32(replications) / float32(n.Manifest.Replications)
			prob := rand.Float32()

			// If the probability of pulling is greater than activation energy
			// then pull the file locally.
			if prob > activationEnergy {
				file.Size, err = n.Node.GetSize(file.ID)
				if err != nil {
					log.Println(err)
					continue
				}
				err = n.Node.Pin(file.ID)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}

		// Update the number of times replicated to the database
		file.Replications = replications

		// Update entry in database
		_, err = n.DB.Update(file)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

// genNumWorkers generates a number of worker processes to optimize efficiency.
// Subtract 2 from the number of cores because of the main thread and the GetAll function.
func genNumWorkers() int {
	if runtime.NumCPU() > 2 {
		return runtime.NumCPU() - 1
	}
	return 1
}
