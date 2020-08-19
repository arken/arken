package ipfs

import (
	"strings"

	"github.com/ipfs/go-ipfs/core/corerepo"
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
	err = corerepo.GarbageCollect(node, ctx)
	if err != nil {
		return err
	}
	return nil
}
