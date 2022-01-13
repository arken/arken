package ipfs

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

func (n *Node) ReportStats(peerID string, protocolID string, msg []byte) (err error) {
	id, err := peer.Decode(peerID)
	if err != nil {
		return err
	}
	_, err = n.node.DHT.FindPeer(n.ctx, id)
	if err != nil {
		return err
	}
	stream, err := n.node.PeerHost.NewStream(n.ctx, id, protocol.ID(protocolID))
	if err != nil {
		return err
	}
	_, err = stream.Write(msg)
	if err != nil {
		return err
	}
	err = stream.Close()
	return err
}
