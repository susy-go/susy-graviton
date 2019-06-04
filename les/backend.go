// Copyleft 2016 The susy-graviton Authors
// This file is part of the susy-graviton library.
//
// The susy-graviton library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The susy-graviton library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MSRCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the susy-graviton library. If not, see <http://www.gnu.org/licenses/>.

// Package les implements the Light Sophon Subprotocol.
package les

import (
	"fmt"
	"sync"
	"time"

	"github.com/susy-go/susy-graviton/accounts"
	"github.com/susy-go/susy-graviton/common"
	"github.com/susy-go/susy-graviton/common/hexutil"
	"github.com/susy-go/susy-graviton/consensus"
	"github.com/susy-go/susy-graviton/core"
	"github.com/susy-go/susy-graviton/core/bloombits"
	"github.com/susy-go/susy-graviton/core/rawdb"
	"github.com/susy-go/susy-graviton/core/types"
	"github.com/susy-go/susy-graviton/sof"
	"github.com/susy-go/susy-graviton/sof/downloader"
	"github.com/susy-go/susy-graviton/sof/filters"
	"github.com/susy-go/susy-graviton/sof/gasprice"
	"github.com/susy-go/susy-graviton/event"
	"github.com/susy-go/susy-graviton/internal/sofapi"
	"github.com/susy-go/susy-graviton/light"
	"github.com/susy-go/susy-graviton/log"
	"github.com/susy-go/susy-graviton/node"
	"github.com/susy-go/susy-graviton/p2p"
	"github.com/susy-go/susy-graviton/p2p/discv5"
	"github.com/susy-go/susy-graviton/params"
	rpc "github.com/susy-go/susy-graviton/rpc"
)

type LightSophon struct {
	lesCommons

	odr         *LesOdr
	relay       *LesTxRelay
	chainConfig *params.ChainConfig
	// Channel for shutting down the service
	shutdownChan chan bool

	// Handlers
	peers      *peerSet
	txPool     *light.TxPool
	blockchain *light.LightChain
	serverPool *serverPool
	reqDist    *requestDistributor
	retriever  *retrieveManager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer

	ApiBackend *LesApiBackend

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	networkId     uint64
	netRPCService *sofapi.PublicNetAPI

	wg sync.WaitGroup
}

func New(ctx *node.ServiceContext, config *sof.Config) (*LightSophon, error) {
	chainDb, err := sof.CreateDB(ctx, config, "lightchaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlockWithOverride(chainDb, config.Genesis, config.ConstantinopleOverride)
	if _, isCompat := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !isCompat {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	peers := newPeerSet()
	quitSync := make(chan struct{})

	lsof := &LightSophon{
		lesCommons: lesCommons{
			chainDb: chainDb,
			config:  config,
			iConfig: light.DefaultClientIndexerConfig,
		},
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		peers:          peers,
		reqDist:        newRequestDistributor(peers, quitSync),
		accountManager: ctx.AccountManager,
		engine:         sof.CreateConsensusEngine(ctx, chainConfig, &config.Sofash, nil, false, chainDb),
		shutdownChan:   make(chan bool),
		networkId:      config.NetworkId,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   sof.NewBloomIndexer(chainDb, params.BloomBitsBlocksClient, params.HelperTrieConfirmations),
	}

	lsof.relay = NewLesTxRelay(peers, lsof.reqDist)
	lsof.serverPool = newServerPool(chainDb, quitSync, &lsof.wg)
	lsof.retriever = newRetrieveManager(peers, lsof.reqDist, lsof.serverPool)

	lsof.odr = NewLesOdr(chainDb, light.DefaultClientIndexerConfig, lsof.retriever)
	lsof.chtIndexer = light.NewChtIndexer(chainDb, lsof.odr, params.CHTFrequencyClient, params.HelperTrieConfirmations)
	lsof.bloomTrieIndexer = light.NewBloomTrieIndexer(chainDb, lsof.odr, params.BloomBitsBlocksClient, params.BloomTrieFrequency)
	lsof.odr.SetIndexers(lsof.chtIndexer, lsof.bloomTrieIndexer, lsof.bloomIndexer)

	// Note: NewLightChain adds the trusted checkpoint so it needs an ODR with
	// indexers already set but not started yet
	if lsof.blockchain, err = light.NewLightChain(lsof.odr, lsof.chainConfig, lsof.engine); err != nil {
		return nil, err
	}
	// Note: AddChildIndexer starts the update process for the child
	lsof.bloomIndexer.AddChildIndexer(lsof.bloomTrieIndexer)
	lsof.chtIndexer.Start(lsof.blockchain)
	lsof.bloomIndexer.Start(lsof.blockchain)

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		lsof.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}

	lsof.txPool = light.NewTxPool(lsof.chainConfig, lsof.blockchain, lsof.relay)
	if lsof.protocolManager, err = NewProtocolManager(lsof.chainConfig, light.DefaultClientIndexerConfig, true, config.NetworkId, lsof.eventMux, lsof.engine, lsof.peers, lsof.blockchain, nil, chainDb, lsof.odr, lsof.relay, lsof.serverPool, quitSync, &lsof.wg); err != nil {
		return nil, err
	}
	lsof.ApiBackend = &LesApiBackend{lsof, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.MinerGasPrice
	}
	lsof.ApiBackend.gpo = gasprice.NewOracle(lsof.ApiBackend, gpoParams)
	return lsof, nil
}

func lesTopic(genesisHash common.Hash, protocolVersion uint) discv5.Topic {
	var name string
	switch protocolVersion {
	case lpv1:
		name = "LES"
	case lpv2:
		name = "LES2"
	default:
		panic(nil)
	}
	return discv5.Topic(name + "@" + common.Bytes2Hex(genesisHash.Bytes()[0:8]))
}

type LightDummyAPI struct{}

// Sophybase is the address that mining rewards will be send to
func (s *LightDummyAPI) Sophybase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Coinbase is the address that mining rewards will be send to (alias for Sophybase)
func (s *LightDummyAPI) Coinbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Hashrate returns the POW hashrate
func (s *LightDummyAPI) Hashrate() hexutil.Uint {
	return 0
}

// Mining returns an indication if this node is currently mining.
func (s *LightDummyAPI) Mining() bool {
	return false
}

// APIs returns the collection of RPC services the sophon package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *LightSophon) APIs() []rpc.API {
	return append(sofapi.GetAPIs(s.ApiBackend), []rpc.API{
		{
			Namespace: "sof",
			Version:   "1.0",
			Service:   &LightDummyAPI{},
			Public:    true,
		}, {
			Namespace: "sof",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "sof",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, true),
			Public:    true,
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *LightSophon) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *LightSophon) BlockChain() *light.LightChain      { return s.blockchain }
func (s *LightSophon) TxPool() *light.TxPool              { return s.txPool }
func (s *LightSophon) Engine() consensus.Engine           { return s.engine }
func (s *LightSophon) LesVersion() int                    { return int(ClientProtocolVersions[0]) }
func (s *LightSophon) Downloader() *downloader.Downloader { return s.protocolManager.downloader }
func (s *LightSophon) EventMux() *event.TypeMux           { return s.eventMux }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *LightSophon) Protocols() []p2p.Protocol {
	return s.makeProtocols(ClientProtocolVersions)
}

// Start implements node.Service, starting all internal goroutines needed by the
// Sophon protocol implementation.
func (s *LightSophon) Start(srvr *p2p.Server) error {
	log.Warn("Light client mode is an experimental feature")
	s.startBloomHandlers(params.BloomBitsBlocksClient)
	s.netRPCService = sofapi.NewPublicNetAPI(srvr, s.networkId)
	// clients are searching for the first advertised protocol in the list
	protocolVersion := AdvertiseProtocolVersions[0]
	s.serverPool.start(srvr, lesTopic(s.blockchain.Genesis().Hash(), protocolVersion))
	s.protocolManager.Start(s.config.LightPeers)
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Sophon protocol.
func (s *LightSophon) Stop() error {
	s.odr.Stop()
	s.bloomIndexer.Close()
	s.chtIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	s.txPool.Stop()
	s.engine.Close()

	s.eventMux.Stop()

	time.Sleep(time.Millisecond * 200)
	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
