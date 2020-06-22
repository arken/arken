package engine

import (
	"github.com/archivalists/arken/ipfs"
)

// CalcThreshold calculates the AtRiskThreshold for the minimum
// number of nodes that should backup a file to ensure data protection.
func CalcThreshold(lightHouseKey string, replicationFactor float32, startNodes int) (threshold int, err error) {
	// ToDo: Change this to 1000 and return 100 after development.
	if startNodes > 1000 {
		return 100, nil
	}
	maxNodes, err := ipfs.FindProvs(lightHouseKey, startNodes)
	if err != nil {
		return -1, err
	}
	if maxNodes == startNodes {
		return CalcThreshold(lightHouseKey, replicationFactor, startNodes*2)
	}
	threshold = int(float32(maxNodes) * replicationFactor)
	if threshold < 5 {
		threshold = 5
	}
	return threshold, nil
}
