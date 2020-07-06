package ipfs

import (
	"path/filepath"

	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

type writeFunc func(peer.AddrInfo) (err error)

// AddBootstrapPeer adds a peer to the IPFS node as an important connection to keep alive.
func AddBootstrapPeer(peerID string) (err error) {
	cid := filepath.Base(peerID)
	addr := filepath.Dir(filepath.Dir(peerID))
	pid, err := peer.Decode(cid)
	if err != nil {
		return err
	}

	address, err := ma.NewMultiaddr(addr)
	if err != nil {
		return err
	}
	peer := peer.AddrInfo{ID: pid, Addrs: []ma.Multiaddr{address}}
	ps.AddPeer(peer)

	return nil
}
