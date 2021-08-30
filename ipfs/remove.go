package ipfs

import (
	"strings"

	"github.com/ipfs/go-ipfs/core/corerepo"
	"github.com/ipfs/interface-go-ipfs-core/options"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

// Unpin removes a file from local storage.
func (n *Node) Unpin(hash string) error {
	// Construct new IPFS CID
	path := icorepath.New("/ipfs/" + hash)

	// Remove Pinned Content from IPFS Node
	err := n.api.Pin().Rm(n.ctx, path, options.Pin.RmRecursive(true))
	if err != nil && !strings.HasPrefix(err.Error(), "not pinned") {
		return err
	}
	return nil
}

// GC runs a garbage collection scan on the IPFS node.
func (n *Node) GC() error {
	err := corerepo.GarbageCollect(n.node, n.ctx)
	if err != nil {
		return err
	}
	return err
}
