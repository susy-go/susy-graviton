// Copyleft 2015 The susy-graviton Authors
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

package sof

import (
	"context"
	"math/big"

	"github.com/susy-go/susy-graviton/accounts"
	"github.com/susy-go/susy-graviton/common"
	"github.com/susy-go/susy-graviton/common/math"
	"github.com/susy-go/susy-graviton/core"
	"github.com/susy-go/susy-graviton/core/bloombits"
	"github.com/susy-go/susy-graviton/core/state"
	"github.com/susy-go/susy-graviton/core/types"
	"github.com/susy-go/susy-graviton/core/vm"
	"github.com/susy-go/susy-graviton/sof/downloader"
	"github.com/susy-go/susy-graviton/sof/gasprice"
	"github.com/susy-go/susy-graviton/sofdb"
	"github.com/susy-go/susy-graviton/event"
	"github.com/susy-go/susy-graviton/params"
	"github.com/susy-go/susy-graviton/rpc"
)

// SofAPIBackend implements sofapi.Backend for full nodes
type SofAPIBackend struct {
	sof *Sophon
	gpo *gasprice.Oracle
}

// ChainConfig returns the active chain configuration.
func (b *SofAPIBackend) ChainConfig() *params.ChainConfig {
	return b.sof.chainConfig
}

func (b *SofAPIBackend) CurrentBlock() *types.Block {
	return b.sof.blockchain.CurrentBlock()
}

func (b *SofAPIBackend) SetHead(number uint64) {
	b.sof.protocolManager.downloader.Cancel()
	b.sof.blockchain.SetHead(number)
}

func (b *SofAPIBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.sof.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.sof.blockchain.CurrentBlock().Header(), nil
	}
	return b.sof.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *SofAPIBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.sof.blockchain.GetHeaderByHash(hash), nil
}

func (b *SofAPIBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.sof.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.sof.blockchain.CurrentBlock(), nil
	}
	return b.sof.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *SofAPIBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.sof.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.sof.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *SofAPIBackend) GetBlock(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.sof.blockchain.GetBlockByHash(hash), nil
}

func (b *SofAPIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	return b.sof.blockchain.GetReceiptsByHash(hash), nil
}

func (b *SofAPIBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	receipts := b.sof.blockchain.GetReceiptsByHash(hash)
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *SofAPIBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.sof.blockchain.GetTdByHash(blockHash)
}

func (b *SofAPIBackend) GetSVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header) (*vm.SVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewSVMContext(msg, header, b.sof.BlockChain(), nil)
	return vm.NewSVM(context, state, b.sof.chainConfig, *b.sof.blockchain.GetVMConfig()), vmError, nil
}

func (b *SofAPIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.sof.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *SofAPIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.sof.BlockChain().SubscribeChainEvent(ch)
}

func (b *SofAPIBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.sof.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *SofAPIBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.sof.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *SofAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.sof.BlockChain().SubscribeLogsEvent(ch)
}

func (b *SofAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.sof.txPool.AddLocal(signedTx)
}

func (b *SofAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.sof.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *SofAPIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.sof.txPool.Get(hash)
}

func (b *SofAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.sof.txPool.State().GetNonce(addr), nil
}

func (b *SofAPIBackend) Stats() (pending int, queued int) {
	return b.sof.txPool.Stats()
}

func (b *SofAPIBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.sof.TxPool().Content()
}

func (b *SofAPIBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.sof.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *SofAPIBackend) Downloader() *downloader.Downloader {
	return b.sof.Downloader()
}

func (b *SofAPIBackend) ProtocolVersion() int {
	return b.sof.SofVersion()
}

func (b *SofAPIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *SofAPIBackend) ChainDb() sofdb.Database {
	return b.sof.ChainDb()
}

func (b *SofAPIBackend) EventMux() *event.TypeMux {
	return b.sof.EventMux()
}

func (b *SofAPIBackend) AccountManager() *accounts.Manager {
	return b.sof.AccountManager()
}

func (b *SofAPIBackend) RPCGasCap() *big.Int {
	return b.sof.config.RPCGasCap
}

func (b *SofAPIBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.sof.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *SofAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.sof.bloomRequests)
	}
}
