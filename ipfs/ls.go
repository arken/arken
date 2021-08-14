package ipfs

import (
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
)

// LsPin returns a list of files IDs that are pinned on the IPFS instance.
func (n *Node) LsPin() (<-chan string, error) {
	// Construct output channel
	output := make(chan string)

	// List all directly pinned files on the node.
	pins, err := n.api.Pin().Ls(n.ctx, options.Pin.Ls.Direct())
	if err != nil {
		return nil, err
	}

	// Construct CID from pin for all files
	go func(pins <-chan iface.Pin) {
		// Loop through pinned files
		for pin := range pins {
			cid := pin.Path().Cid().String()
			output <- cid
		}
		close(output)
	}(pins)

	// Return output channel
	return output, nil
}
