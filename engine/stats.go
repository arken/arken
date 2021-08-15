package engine

import (
	"encoding/json"
	"log"

	"github.com/dustin/go-humanize"
)

type report struct {
	Email      string  `json:"email"`
	TotalSpace float64 `json:"total_space"`
	UsedSpace  float64 `json:"used_space"`
}

func (n *Node) ReportStats() {
	gb, err := humanize.ParseBytes("1GB")
	if err != nil {
		log.Println(err)
		return
	}

	totalSpace, err := humanize.ParseBytes(n.Cfg.Storage.Limit)
	if err != nil {
		log.Println(err)
		return
	}

	usedSpace, err := n.Node.RepoSize()
	if err != nil {
		log.Println(err)
		return
	}

	report := report{
		Email:      n.Cfg.Stats.Email,
		UsedSpace:  float64(usedSpace) / float64(gb),
		TotalSpace: float64(totalSpace) / float64(gb),
	}

	bytes, err := json.Marshal(report)
	if err != nil {
		log.Println(err)
		return
	}

	err = n.Node.ReportStats(n.Manifest.StatsNode, "/arkstat/0.0.1", bytes)
	if err != nil {
		log.Println(err)
		return
	}
}
