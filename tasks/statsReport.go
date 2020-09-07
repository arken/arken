package tasks

import (
	"log"
	"time"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/stats"
)

// StatsReporting periodically reports node statistics for
// configured keysets.
func StatsReporting() {
	// If allowed report the stats to the keyset stats server.
	err := stats.Report(config.Keysets)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(1 * time.Hour)
}
