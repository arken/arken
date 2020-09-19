package tasks

import (
	"time"

	"github.com/arkenproject/arken/engine"
)

// VerifyLocal runs a check hourly to verify locally pinned files are
// still present on the system.
func VerifyLocal() {
	for {
		engine.VerifyLocal()
		time.Sleep(24 * time.Hour)
	}
}
