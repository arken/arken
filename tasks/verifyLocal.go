package tasks

import (
	"time"

	"github.com/arkenproject/arken/engine"
)

// VerifyLocal runs a weekly check to verify locally pinned files are
// still present on the system.
func VerifyLocal() {
	for {
		time.Sleep(1 * time.Hour)
		engine.VerifyLocal()
		time.Sleep(7 * 23 * time.Hour)
	}
}
