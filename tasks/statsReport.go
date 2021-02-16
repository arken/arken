package tasks

import (
	"log"
	"time"

	"github.com/arken/arken/config"
	"github.com/arken/arken/stats"
)

// StatsReporting periodically reports node statistics for
// configured keysets.
func StatsReporting() {
	for {
		// If allowed report the stats to the keyset stats server.
		err := stats.Report(config.Keysets)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(1 * time.Hour)
	}
}
