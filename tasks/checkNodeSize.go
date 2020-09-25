package tasks

import (
	"fmt"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
	"github.com/arkenproject/arken/engine"
	"github.com/arkenproject/arken/ipfs"
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
