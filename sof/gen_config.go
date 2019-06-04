// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package sof

import (
	"math/big"
	"time"

	"github.com/susy-go/susy-graviton/common"
	"github.com/susy-go/susy-graviton/consensus/sofash"
	"github.com/susy-go/susy-graviton/core"
	"github.com/susy-go/susy-graviton/sof/downloader"
	"github.com/susy-go/susy-graviton/sof/gasprice"
	"github.com/susy-go/susy-graviton/miner"
)

// MarshalTOML marshals as TOML.
func (c Config) MarshalTOML() (interface{}, error) {
	type Config struct {
		Genesis                 *core.Genesis `toml:",omitempty"`
		NetworkId               uint64
		SyncMode                downloader.SyncMode
		NoPruning               bool
		NoPrefetch              bool
		Whitelist               map[uint64]common.Hash `toml:"-"`
		LightServ               int                    `toml:",omitempty"`
		LightBandwidthIn        int                    `toml:",omitempty"`
		LightBandwidthOut       int                    `toml:",omitempty"`
		LightPeers              int                    `toml:",omitempty"`
		OnlyAnnounce            bool
		ULC                     *ULCConfig `toml:",omitempty"`
		SkipBcVersionCheck      bool       `toml:"-"`
		DatabaseHandles         int        `toml:"-"`
		DatabaseCache           int
		TrieCleanCache          int
		TrieDirtyCache          int
		TrieTimeout             time.Duration
		Miner                   miner.Config
		Sofash                  sofash.Config
		TxPool                  core.TxPoolConfig
		GPO                     gasprice.Config
		EnablePreimageRecording bool
		DocRoot                 string `toml:"-"`
		EWASMInterpreter        string
		SVMInterpreter          string
		ConstantinopleOverride  *big.Int
		RPCGasCap               *big.Int `toml:",omitempty"`
	}
	var enc Config
	enc.Genesis = c.Genesis
	enc.NetworkId = c.NetworkId
	enc.SyncMode = c.SyncMode
	enc.NoPruning = c.NoPruning
	enc.NoPrefetch = c.NoPrefetch
	enc.Whitelist = c.Whitelist
	enc.LightServ = c.LightServ
	enc.LightBandwidthIn = c.LightBandwidthIn
	enc.LightBandwidthOut = c.LightBandwidthOut
	enc.LightPeers = c.LightPeers
	enc.OnlyAnnounce = c.OnlyAnnounce
	enc.ULC = c.ULC
	enc.SkipBcVersionCheck = c.SkipBcVersionCheck
	enc.DatabaseHandles = c.DatabaseHandles
	enc.DatabaseCache = c.DatabaseCache
	enc.TrieCleanCache = c.TrieCleanCache
	enc.TrieDirtyCache = c.TrieDirtyCache
	enc.TrieTimeout = c.TrieTimeout
	enc.Miner = c.Miner
	enc.Sofash = c.Sofash
	enc.TxPool = c.TxPool
	enc.GPO = c.GPO
	enc.EnablePreimageRecording = c.EnablePreimageRecording
	enc.DocRoot = c.DocRoot
	enc.EWASMInterpreter = c.EWASMInterpreter
	enc.SVMInterpreter = c.SVMInterpreter
	enc.ConstantinopleOverride = c.ConstantinopleOverride
	enc.RPCGasCap = c.RPCGasCap
	return &enc, nil
}

// UnmarshalTOML unmarshals from TOML.
func (c *Config) UnmarshalTOML(unmarshal func(interface{}) error) error {
	type Config struct {
		Genesis                 *core.Genesis `toml:",omitempty"`
		NetworkId               *uint64
		SyncMode                *downloader.SyncMode
		NoPruning               *bool
		NoPrefetch              *bool
		Whitelist               map[uint64]common.Hash `toml:"-"`
		LightServ               *int                   `toml:",omitempty"`
		LightBandwidthIn        *int                   `toml:",omitempty"`
		LightBandwidthOut       *int                   `toml:",omitempty"`
		LightPeers              *int                   `toml:",omitempty"`
		OnlyAnnounce            *bool
		ULC                     *ULCConfig `toml:",omitempty"`
		SkipBcVersionCheck      *bool      `toml:"-"`
		DatabaseHandles         *int       `toml:"-"`
		DatabaseCache           *int
		TrieCleanCache          *int
		TrieDirtyCache          *int
		TrieTimeout             *time.Duration
		Miner                   *miner.Config
		Sofash                  *sofash.Config
		TxPool                  *core.TxPoolConfig
		GPO                     *gasprice.Config
		EnablePreimageRecording *bool
		DocRoot                 *string `toml:"-"`
		EWASMInterpreter        *string
		SVMInterpreter          *string
		ConstantinopleOverride  *big.Int
		RPCGasCap               *big.Int `toml:",omitempty"`
	}
	var dec Config
	if err := unmarshal(&dec); err != nil {
		return err
	}
	if dec.Genesis != nil {
		c.Genesis = dec.Genesis
	}
	if dec.NetworkId != nil {
		c.NetworkId = *dec.NetworkId
	}
	if dec.SyncMode != nil {
		c.SyncMode = *dec.SyncMode
	}
	if dec.NoPruning != nil {
		c.NoPruning = *dec.NoPruning
	}
	if dec.NoPrefetch != nil {
		c.NoPrefetch = *dec.NoPrefetch
	}
	if dec.Whitelist != nil {
		c.Whitelist = dec.Whitelist
	}
	if dec.LightServ != nil {
		c.LightServ = *dec.LightServ
	}
	if dec.LightBandwidthIn != nil {
		c.LightBandwidthIn = *dec.LightBandwidthIn
	}
	if dec.LightBandwidthOut != nil {
		c.LightBandwidthOut = *dec.LightBandwidthOut
	}
	if dec.LightPeers != nil {
		c.LightPeers = *dec.LightPeers
	}
	if dec.OnlyAnnounce != nil {
		c.OnlyAnnounce = *dec.OnlyAnnounce
	}
	if dec.ULC != nil {
		c.ULC = dec.ULC
	}
	if dec.SkipBcVersionCheck != nil {
		c.SkipBcVersionCheck = *dec.SkipBcVersionCheck
	}
	if dec.DatabaseHandles != nil {
		c.DatabaseHandles = *dec.DatabaseHandles
	}
	if dec.DatabaseCache != nil {
		c.DatabaseCache = *dec.DatabaseCache
	}
	if dec.TrieCleanCache != nil {
		c.TrieCleanCache = *dec.TrieCleanCache
	}
	if dec.TrieDirtyCache != nil {
		c.TrieDirtyCache = *dec.TrieDirtyCache
	}
	if dec.TrieTimeout != nil {
		c.TrieTimeout = *dec.TrieTimeout
	}
	if dec.Miner != nil {
		c.Miner = *dec.Miner
	}
	if dec.Sofash != nil {
		c.Sofash = *dec.Sofash
	}
	if dec.TxPool != nil {
		c.TxPool = *dec.TxPool
	}
	if dec.GPO != nil {
		c.GPO = *dec.GPO
	}
	if dec.EnablePreimageRecording != nil {
		c.EnablePreimageRecording = *dec.EnablePreimageRecording
	}
	if dec.DocRoot != nil {
		c.DocRoot = *dec.DocRoot
	}
	if dec.EWASMInterpreter != nil {
		c.EWASMInterpreter = *dec.EWASMInterpreter
	}
	if dec.SVMInterpreter != nil {
		c.SVMInterpreter = *dec.SVMInterpreter
	}
	if dec.ConstantinopleOverride != nil {
		c.ConstantinopleOverride = dec.ConstantinopleOverride
	}
	if dec.RPCGasCap != nil {
		c.RPCGasCap = dec.RPCGasCap
	}
	return nil
}
