package engine

import (
	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
)

// CheckNetUsage returns if the total number of files downloaded in the last
// 31 days equals or is greater than the limit.
func CheckNetUsage() (hit bool, err error) {
	db, err := database.Open(config.Global.Database.Path)
	if err != nil {
		return true, err
	}

	sum, err := database.TransactionSum(db, "added")
	if err != nil {
		return true, err
	}
	return sum > config.ParseWellFormedPoolSize(config.Global.General.NetworkLimit), err
}
