package tasks

import (
	"fmt"
	"log"
	"time"

	"github.com/arken/arken/engine"
)

func checkNetworkLimit() {
	var (
		err error
	)

	limit, err := engine.CheckNetUsage()
	if err != nil {
		log.Fatal(err)
	}
	if limit && !engine.NetworkLimit {
		fmt.Printf("\n[Cancelling File Sync due to Network Limit Hit]\n")
	}
	if !limit && engine.NetworkLimit {
		fmt.Printf("\n[Network Limit Below Threshold, Resuming File Sync Operations]\n")
	}

	engine.NetworkLimit = limit

	time.Sleep(1 * time.Hour)
}
