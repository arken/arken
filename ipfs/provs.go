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
	contxt, cancl := context.WithTimeout(ctx, 5*time.Second)
	path := icorepath.New("/ipfs/" + hash)
	output, err := ipfs.Dht().FindProviders(contxt, path, func(input *options.DhtFindProvidersSettings) error {
		input.NumProviders = maxPeers + 10
		return nil
	})
	if err != nil {
		cancl()
		return -1, err
	}

	count := 0
	for range output {
		_ = <-output
		count++
	}
	cancl()
	return count, nil
}
