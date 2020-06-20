package ipfs

import (
	"github.com/ipfs/interface-go-ipfs-core/options"

	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

// FindProvs queries the IPFS network for the number of
// providers hosting a given file
func FindProvs(hash string) (replications int, err error) {
	path := icorepath.New("/ipfs/" + hash)
	output, err := ipfs.Dht().FindProviders(ctx, path, func(input *options.DhtFindProvidersSettings) error {
		input.NumProviders = AtRiskThreshhold
		return nil
	})
	if err != nil {
		return -1, err
	}

	count := 0
	for range output {
		count++
	}

	return count, nil
}
