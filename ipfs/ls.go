package ipfs

import (
	"github.com/arken/arken/database"
	"github.com/ipfs/interface-go-ipfs-core/options"
)

// LsPin returns a list of File that are pinned on the IPFS instance.
func (n *Node) LsPin(output chan<- database.File, errors chan<- error) {
	defer close(output)

	// List all directly pinned files on the node.
	pins, err := n.api.Pin().Ls(n.ctx, options.Pin.Ls.Direct())
	if err != nil {
		errors <- err
	}

	// Loop through pinned files
	for pin := range pins {
		cid := pin.Path().Cid().String()
		size, err := n.GetSize(cid)
		if err != nil {
			continue
		}

		fileTemplate := database.File{
			ID:           cid,
			Size:         int64(size),
			Status:       "local",
			Replications: -1,
		}

		output <- fileTemplate
	}
	errors <- nil
}
