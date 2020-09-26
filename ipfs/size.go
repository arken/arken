package ipfs

import (
	"context"
	"time"

	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

// GetSize returns the size of the specified file.
func GetSize(hash string) (size int, err error) {
	path := icorepath.New("/ipfs/" + hash)
	contxt, cancl := context.WithTimeout(ctx, 20*time.Second)

	obj, err := ipfs.Object().Stat(contxt, path)
	if err != nil {
		cancl()
		return -1, err
	}
	err = contxt.Err()
	if err != nil {
		cancl()
		return -1, err
	}
	cancl()
	return obj.CumulativeSize, err
}
