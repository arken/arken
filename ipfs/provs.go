package ipfs

import (
	"context"
	"time"

	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"

	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

// FindProvs queries the IPFS network for the number of
// providers hosting a given file
func FindProvs(hash string, maxPeers int) (replications int, err error) {
	path := icorepath.New("/ipfs/" + hash)
	contxt, cancl := context.WithTimeout(ctx, 5*time.Second)
	contxt, events := routing.RegisterForQueryEvents(contxt)

	output, err := ipfs.Dht().FindProviders(contxt, path, func(input *options.DhtFindProvidersSettings) error {
		input.NumProviders = maxPeers + 15
		return nil
	})
	if err != nil {
		cancl()
		return -1, err
	}
	go func() {
		defer cancl()
		for p := range output {
			np := p
			routing.PublishQueryEvent(ctx, &routing.QueryEvent{
				Type:      routing.Provider,
				Responses: []*peer.AddrInfo{&np},
			})
		}
	}()

	set := make(map[string]bool)
	for node := range events {
		set[node.ID.String()] = true
	}

	return len(set), nil
}
