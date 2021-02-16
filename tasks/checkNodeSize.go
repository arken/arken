package tasks

import (
	"fmt"

	"github.com/arken/arken/config"
	"github.com/arken/arken/database"
	"github.com/arken/arken/engine"
	"github.com/arken/arken/ipfs"
)

func checkNodeSize(output chan<- database.FileKey) error {
	repoSize, err := ipfs.GetRepoSize()
	if err != nil {
		return err
	}

	poolSize := config.ParseWellFormedPoolSize(config.Global.General.PoolSize)
	if err != nil {
		return err
	}

	if repoSize > poolSize {
		fmt.Println("[Reducing Node Usage on Disk]")
		request := repoSize - poolSize
		err = engine.MakeSpace(int64(request), output, true)
		if err != nil {
			return err
		}
	}
	return nil
}
