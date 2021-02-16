package engine

import (
	"github.com/arken/arken/ipfs"
)

// CalcThreshold calculates the AtRiskThreshold for the minimum
// number of nodes that should backup a file to ensure data protection.
func CalcThreshold(lightHouseKey string, replications int, startNodes int) (threshold int, err error) {
	if startNodes > 1000 {
		return 100, nil
	}
	if replications < 0 {
		maxNodes, err := ipfs.FindProvs(lightHouseKey, startNodes)
		if err != nil {
			return -1, err
		}
		if maxNodes == startNodes {
			return CalcThreshold(lightHouseKey, replications, startNodes*2)
		}
		return maxNodes, nil
	}
	return replications, nil
}
