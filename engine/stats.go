package engine

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/disk"
)

type report struct {
	Email      string  `json:"email"`
	TotalSpace float64 `json:"total_space"`
	UsedSpace  float64 `json:"used_space"`
}

func (n *Node) ReportStats() {
	fmt.Printf("Sending Stats to Cluster...\n")

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

	dstat, err := disk.Usage(n.Cfg.Manifest.Path)
	if err != nil {
		log.Println(err)
		return
	}

	// Don't allow the user to contribute more space than they have.
	if dstat.Total < totalSpace {
		totalSpace = dstat.Total
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
	fmt.Printf("Stats sent successfully\n")
}
