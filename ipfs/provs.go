package ipfs

import (
	"context"
	"time"

	"github.com/ipfs/interface-go-ipfs-core/options"

	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

// FindProvs queries the IPFS network for the number of
// providers hosting a given file
func (n *Node) FindProvs(hash string, maxPeers int) (replications int, err error) {
	// Construct IPFS CID
	path := icorepath.New("/ipfs/" + hash)

	// Create a new context
	contxt, cancel := context.WithTimeout(n.ctx, 5*time.Second)
	defer cancel()

	// Lookup how many other nodes are hosting a file.
	output, err := n.api.Dht().FindProviders(contxt, path, options.Dht.NumProviders(maxPeers+15))
	if err != nil {
		return -1, err
	}

	// Iterate through resulting responses and add them up.
	for range output {
		replications++
	}

	return replications, nil
}
