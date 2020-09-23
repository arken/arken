package ipfs

import (
	"context"
	"time"

	"github.com/ipfs/interface-go-ipfs-core/options"

	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

// FindProvs queries the IPFS network for the number of
// providers hosting a given file
func FindProvs(hash string, maxPeers int) (replications int, err error) {
	path := icorepath.New("/ipfs/" + hash)
	contxt, cancl := context.WithTimeout(ctx, 5*time.Second)

	output, err := ipfs.Dht().FindProviders(contxt, path, func(input *options.DhtFindProvidersSettings) error {
		input.NumProviders = maxPeers + 15
		return nil
	})
	if err != nil {
		cancl()
		return -1, err
	}

	for range output {
		replications++
	}

	cancl()
	return replications, nil
}
