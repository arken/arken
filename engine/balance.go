package engine

import (
	"database/sql"
	"fmt"
	"log"
)

func (n *Node) Rebalance() {
	// Boot adder subsystem
	toAdder, err := n.FileAdder()
	if err != nil {
		log.Println(err)
	}
	defer close(toAdder)

	for i := 0; ; i++ {
		// Get all remote files from DB
		files, err := n.DB.GetAll("remote", 100, i)
		if err != nil {
			break
		}

		// Loop through files in list
		for _, file := range files {
			toAdder <- file
		}
	}
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
	}
}
