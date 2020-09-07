package engine

import (
	"os"
	"path/filepath"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
)

// CheckNetUsage returns if the total number of files downloaded in the last
// 31 days equals or is greater than the limit.
func CheckNetUsage() (hit bool, err error) {
	copyName := filepath.Join(filepath.Dir(config.Global.Database.Path), "NetCheck.db")
	database.Copy(config.Global.Database.Path, copyName)
	defer os.Remove(copyName)

	db, err := database.Open(copyName)
	if err != nil {
		return true, err
	}

	sum, err := database.TransactionSum(db, "added")
	if err != nil {
		return true, err
	}
	return sum > config.ParseWellFormedPoolSize(config.Global.General.NetworkLimit), err
}
