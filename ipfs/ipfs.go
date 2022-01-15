package ipfs

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	ipfsConfig "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi" // This package is needed so that all the preloaded plugins are loaded automatically.
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	migrate "github.com/ipfs/go-ipfs/repo/fsrepo/migrations"
	icore "github.com/ipfs/interface-go-ipfs-core"
)

type NodeConfArgs struct {
	SwarmKey       string
	StorageMax     string
	BootstrapPeers []string
}

type Node struct {
	api    icore.CoreAPI
	ctx    context.Context
	cancel context.CancelFunc
	node   *core.IpfsNode
}

// CreateNode creates an IPFS node and returns its coreAPI
func CreateNode(repoPath string, args NodeConfArgs) (node *Node, err error) {
	// Setup IPFS plugins
	if err := setupPlugins(repoPath); err != nil {
		return nil, err
	}
	// Initialize node structure
	node = &Node{}
	// Create IPFS node
	node.ctx, node.cancel = context.WithCancel(context.Background())
	// Create Swarm Key File
	if args.SwarmKey != "" {
		err = createSwarmKey(repoPath, args.SwarmKey)
		if err != nil {
			return nil, err
		}
	}
	// Open the repo
	fs, err := openFs(node.ctx, repoPath)
	if err != nil {
		err = createFs(
			node.ctx,
			repoPath,
			args.StorageMax,
			args.BootstrapPeers,
		)
		if err != nil {
			return nil, err
		}
		fs, err = openFs(node.ctx, repoPath)
		if err != nil {
			return nil, err
		}
	}
	// Construct the node
	nodeOptions := &core.BuildCfg{
		Permanent: true,
		Online:    true,
		Routing:   libp2p.DHTClientOption,
		Repo:      fs,
	}
	node.node, err = core.NewNode(node.ctx, nodeOptions)
	if err != nil {
		return nil, err
	}
	node.node.IsDaemon = true

	// Bootstrap the DHT table for peers.
	err = node.node.DHT.Bootstrap(node.ctx)
	if err != nil {
		return nil, err
	}

	// Attach the Core API to the constructed node
	node.api, err = coreapi.NewCoreAPI(node.node)
	return node, err

}

func openFs(ctx context.Context, repoPath string) (result repo.Repo, err error) {
	result, err = fsrepo.Open(repoPath)
	if err != nil && err == fsrepo.ErrNeedMigration {
		err = os.Setenv("IPFS_PATH", repoPath)
		if err != nil {
			return nil, err
		}
		err = migrate.RunMigration(ctx, migrate.NewHttpFetcher("", "", "ipfs", 0), fsrepo.RepoVersion, repoPath, false)
		if err != nil {
			return nil, err
		}
		result, err = fsrepo.Open(repoPath)
	}
	return result, err
}

// createFs builds the IPFS configuration repository.
func createFs(ctx context.Context, path string, storageMax string, bootstrapPeers []string) (err error) {
	// Check if directory to configuration exists
	if _, err = os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	// Create a ipfsConfig with default options and a 2048 bit key
	cfg, err := ipfsConfig.Init(ioutil.Discard, 2048)
	if err != nil {
		return err
	}
	// Set default ipfsConfig values
	cfg.Datastore.StorageMax = storageMax
	cfg.Reprovider.Strategy = "roots"
	cfg.Routing.Type = "dhtclient"
	cfg.Bootstrap = bootstrapPeers
	cfg.Swarm.ConnMgr.LowWater = 20
	cfg.Swarm.ConnMgr.HighWater = 40
	cfg.Swarm.ConnMgr.GracePeriod = time.Minute.String()

	// Create the repo with the ipfsConfig
	err = fsrepo.Init(path, cfg)
	if err != nil {
		return fmt.Errorf("failed to init node: %s", err)
	}
	return nil
}

func createSwarmKey(path string, key string) (err error) {
	keyPath := filepath.Join(path, "swarm.key")
	// Check if directory to configuration exists
	if _, err = os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	if _, err = os.Stat(keyPath); os.IsNotExist(err) {
		var file *os.File
		file, err = os.Create(keyPath)
		if err != nil {
			return err
		}
		_, err = file.WriteString("/key/swarm/psk/1.0.0/\n/base16/\n" + key)
	}
	return err
}

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}
