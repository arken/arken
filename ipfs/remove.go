package ipfs

import (
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
	return err
}
