package tasks

import (
	"fmt"
	"log"
	"time"

	"github.com/arkenproject/arken/engine"
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
		fmt.Printf("\n[Cancelling Rebalance due to Network Limit Hit]\n")
	}
	if !limit && engine.NetworkLimit {
		fmt.Printf("\n[Network Limit Below Threshold, Resuming Balance Operations]\n")
	}

	engine.NetworkLimit = limit

	time.Sleep(1 * time.Hour)
}
