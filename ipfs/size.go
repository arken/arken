package ipfs

import (
	"context"
	"time"

	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

// GetSize returns the size of the specified file.
func (n *Node) GetSize(hash string) (size int64, err error) {
	// Construct new IPFS CID
	path := icorepath.New("/ipfs/" + hash)

	// Create new context with a timeout
	contxt, cancel := context.WithTimeout(n.ctx, 20*time.Second)
	defer cancel()

	// Get stats on an object.
	obj, err := n.api.Object().Stat(contxt, path)
	if err != nil {
		return -1, err
	}

	// Check for errors within the context.
	err = contxt.Err()
	if err != nil {
		return -1, err
	}
	return int64(obj.CumulativeSize), err
}
