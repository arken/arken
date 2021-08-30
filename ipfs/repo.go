package ipfs

func (n *Node) RepoSize() (int64, error) {
	out, err := n.node.Repo.GetStorageUsage()
	return int64(out), err
}
