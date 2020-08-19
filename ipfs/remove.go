package ipfs

import (
	"context"
	"fmt"
	"strings"

	"github.com/ipfs/go-ipfs/gc"
	"github.com/ipfs/interface-go-ipfs-core/options"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

// Unpin removes a file from local storage.
func Unpin(hash string) (err error) {
	path := icorepath.New("/ipfs/" + hash)

	err = ipfs.Pin().Rm(ctx, path, func(input *options.PinRmSettings) error {
		input.Recursive = true
		return nil
	})
	if err != nil && !strings.HasPrefix(err.Error(), "not pinned") {
		return err
	}
	var gcout <-chan gc.Result
	go func() {
		gcout = gc.GC(context.Background(), node.Blockstore, node.Repo.Datastore(), node.Pinning, nil)
	}()
	for cid := range gcout {
		fmt.Printf("Garbage Collected: %s\n", cid)
	}
	return nil
}
