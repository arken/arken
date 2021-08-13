package ipfs

import (
	"github.com/ipfs/interface-go-ipfs-core/options"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

// Pin a file to local storage.
func (n *Node) Pin(hash string) error {
	// Construct IPFS CID
	path := icorepath.New("/ipfs/" + hash)

	// Pin file to local storage within IPFS
	err := n.api.Pin().Add(n.ctx, path, options.Pin.Recursive(true))
	return err
}
