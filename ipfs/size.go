package ipfs

import (
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

// GetSize returns the size of the specified file.
func GetSize(hash string) (size int, err error) {
	path := icorepath.New("/ipfs/" + hash)

	obj, err := ipfs.Object().Stat(ctx, path)
	return obj.CumulativeSize, err
}
