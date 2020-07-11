package ipfs

import (
	"github.com/ipfs/interface-go-ipfs-core/options"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

// Pin a file to local storage.
func Pin(hash string) (err error) {
	path := icorepath.New("/ipfs/" + hash)

	err = ipfs.Pin().Add(ctx, path, func(input *options.PinAddSettings) error {
		input.Recursive = true
		return nil
	})
	return err
}
