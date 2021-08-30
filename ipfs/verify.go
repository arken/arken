package ipfs

// Verify checks the local file store in order to look for corruption.
func (n *Node) Verify() (<-chan string, error) {
	// Initialize out channel of pins
	out := make(chan string, 50)

	// Return a list of pins on the node.
	pins, err := n.api.Pin().Verify(n.ctx)
	if err != nil {
		return nil, err
	}

	// Loop through pins to find broken ones
	go func() {
		defer close(out)
		for pin := range pins {
			if !pin.Ok() {
				for _, bad := range pin.BadNodes() {
					out <- bad.Path().Cid().String()
				}
			}
		}
	}()

	return out, nil
}
