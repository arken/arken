package ipfs

import (
	"log"
	"strings"

	"github.com/arken/arken/database"
	"github.com/ipfs/go-ipfs/core/corerepo"
	"github.com/ipfs/interface-go-ipfs-core/options"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

// Unpin removes a file from local storage.
func Unpin(hash string) (err error) {
	path := icorepath.New("/ipfs/" + hash)

	err = ipfs.Pin().Rm(ctx, path, func(input *options.PinRmSettings) error {
		input.Recursive = true
		return nil
	})
	if err != nil && !strings.HasPrefix(err.Error(), "not pinned") {
		return err
	}
	return nil
}

// LsPin returns a list of FileKeys that are pinned on the IPFS instance.
func LsPin(output chan<- database.FileKey) {
	pins, err := ipfs.Pin().Ls(ctx, func(input *options.PinLsSettings) error {
		input.Type = "recursive"
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	for pin := range pins {
		cid := pin.Path().Cid().String()
		size, err := GetSize(cid)
		if err != nil {
			log.Println(err)
			continue
		}

		fileTemplate := database.FileKey{
			ID:           cid,
			Size:         size,
			Status:       "local",
			Replications: -1,
		}

		output <- fileTemplate
	}

	close(output)
}

// GC runs a garbage collection scan on the IPFS node.
func GC() (err error) {
	err = corerepo.GarbageCollect(node, ctx)
	if err != nil {
		return err
	}
	return err
}
