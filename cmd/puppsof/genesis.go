// Copyleft 2017 The susy-graviton Authors
// This file is part of susy-graviton.
//
// susy-graviton is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// susy-graviton is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MSRCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with susy-graviton. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/binary"
	"errors"
	"math"
	"math/big"
	"strings"

	"github.com/susy-go/susy-graviton/common"
	"github.com/susy-go/susy-graviton/common/hexutil"
	math2 "github.com/susy-go/susy-graviton/common/math"
	"github.com/susy-go/susy-graviton/consensus/sofash"
	"github.com/susy-go/susy-graviton/core"
	"github.com/susy-go/susy-graviton/params"
)

// alsofGenesisSpec represents the genesis specification format used by the
// C++ Sophon implementation.
type alsofGenesisSpec struct {
	SealEngine string `json:"sealEngine"`
	Params     struct {
		AccountStartNonce       math2.HexOrDecimal64   `json:"accountStartNonce"`
		MaximumExtraDataSize    hexutil.Uint64         `json:"maximumExtraDataSize"`
		HomesteadForkBlock      hexutil.Uint64         `json:"homesteadForkBlock"`
		DaoHardforkBlock        math2.HexOrDecimal64   `json:"daoHardforkBlock"`
		SIP150ForkBlock         hexutil.Uint64         `json:"SIP150ForkBlock"`
		SIP158ForkBlock         hexutil.Uint64         `json:"SIP158ForkBlock"`
		ByzantiumForkBlock      hexutil.Uint64         `json:"byzantiumForkBlock"`
		ConstantinopleForkBlock hexutil.Uint64         `json:"constantinopleForkBlock"`
		MinGasLimit             hexutil.Uint64         `json:"minGasLimit"`
		MaxGasLimit             hexutil.Uint64         `json:"maxGasLimit"`
		TieBreakingGas          bool                   `json:"tieBreakingGas"`
		GasLimitBoundDivisor    math2.HexOrDecimal64   `json:"gasLimitBoundDivisor"`
		MinimumDifficulty       *hexutil.Big           `json:"minimumDifficulty"`
		DifficultyBoundDivisor  *math2.HexOrDecimal256 `json:"difficultyBoundDivisor"`
		DurationLimit           *math2.HexOrDecimal256 `json:"durationLimit"`
		BlockReward             *hexutil.Big           `json:"blockReward"`
		NetworkID               hexutil.Uint64         `json:"networkID"`
		ChainID                 hexutil.Uint64         `json:"chainID"`
		AllowFutureBlocks       bool                   `json:"allowFutureBlocks"`
	} `json:"params"`

	Genesis struct {
		Nonce      hexutil.Bytes  `json:"nonce"`
		Difficulty *hexutil.Big   `json:"difficulty"`
		MixHash    common.Hash    `json:"mixHash"`
		Author     common.Address `json:"author"`
		Timestamp  hexutil.Uint64 `json:"timestamp"`
		ParentHash common.Hash    `json:"parentHash"`
		ExtraData  hexutil.Bytes  `json:"extraData"`
		GasLimit   hexutil.Uint64 `json:"gasLimit"`
	} `json:"genesis"`

	Accounts map[common.UnprefixedAddress]*alsofGenesisSpecAccount `json:"accounts"`
}

// alsofGenesisSpecAccount is the prefunded genesis account and/or precompiled
// contract definition.
type alsofGenesisSpecAccount struct {
	Balance     *math2.HexOrDecimal256   `json:"balance"`
	Nonce       uint64                   `json:"nonce,omitempty"`
	Precompiled *alsofGenesisSpecBuiltin `json:"precompiled,omitempty"`
}

// alsofGenesisSpecBuiltin is the precompiled contract definition.
type alsofGenesisSpecBuiltin struct {
	Name          string                         `json:"name,omitempty"`
	StartingBlock hexutil.Uint64                 `json:"startingBlock,omitempty"`
	Linear        *alsofGenesisSpecLinearPricing `json:"linear,omitempty"`
}

type alsofGenesisSpecLinearPricing struct {
	Base uint64 `json:"base"`
	Word uint64 `json:"word"`
}

// newAlsofGenesisSpec converts a susy-graviton genesis block into a Alsof-specific
// chain specification format.
func newAlsofGenesisSpec(network string, genesis *core.Genesis) (*alsofGenesisSpec, error) {
	// Only sofash is currently supported between susy-graviton and alsof
	if genesis.Config.Sofash == nil {
		return nil, errors.New("unsupported consensus engine")
	}
	// Reconstruct the chain spec in Alsof format
	spec := &alsofGenesisSpec{
		SealEngine: "Sofash",
	}
	// Some defaults
	spec.Params.AccountStartNonce = 0
	spec.Params.TieBreakingGas = false
	spec.Params.AllowFutureBlocks = false
	spec.Params.DaoHardforkBlock = 0

	spec.Params.HomesteadForkBlock = (hexutil.Uint64)(genesis.Config.HomesteadBlock.Uint64())
	spec.Params.SIP150ForkBlock = (hexutil.Uint64)(genesis.Config.SIP150Block.Uint64())
	spec.Params.SIP158ForkBlock = (hexutil.Uint64)(genesis.Config.SIP158Block.Uint64())

	// Byzantium
	if num := genesis.Config.ByzantiumBlock; num != nil {
		spec.setByzantium(num)
	}
	// Constantinople
	if num := genesis.Config.ConstantinopleBlock; num != nil {
		spec.setConstantinople(num)
	}

	spec.Params.NetworkID = (hexutil.Uint64)(genesis.Config.ChainID.Uint64())
	spec.Params.ChainID = (hexutil.Uint64)(genesis.Config.ChainID.Uint64())
	spec.Params.MaximumExtraDataSize = (hexutil.Uint64)(params.MaximumExtraDataSize)
	spec.Params.MinGasLimit = (hexutil.Uint64)(params.MinGasLimit)
	spec.Params.MaxGasLimit = (hexutil.Uint64)(math.MaxInt64)
	spec.Params.MinimumDifficulty = (*hexutil.Big)(params.MinimumDifficulty)
	spec.Params.DifficultyBoundDivisor = (*math2.HexOrDecimal256)(params.DifficultyBoundDivisor)
	spec.Params.GasLimitBoundDivisor = (math2.HexOrDecimal64)(params.GasLimitBoundDivisor)
	spec.Params.DurationLimit = (*math2.HexOrDecimal256)(params.DurationLimit)
	spec.Params.BlockReward = (*hexutil.Big)(sofash.FrontierBlockReward)

	spec.Genesis.Nonce = (hexutil.Bytes)(make([]byte, 8))
	binary.LittleEndian.PutUint64(spec.Genesis.Nonce[:], genesis.Nonce)

	spec.Genesis.MixHash = genesis.Mixhash
	spec.Genesis.Difficulty = (*hexutil.Big)(genesis.Difficulty)
	spec.Genesis.Author = genesis.Coinbase
	spec.Genesis.Timestamp = (hexutil.Uint64)(genesis.Timestamp)
	spec.Genesis.ParentHash = genesis.ParentHash
	spec.Genesis.ExtraData = (hexutil.Bytes)(genesis.ExtraData)
	spec.Genesis.GasLimit = (hexutil.Uint64)(genesis.GasLimit)

	for address, account := range genesis.Alloc {
		spec.setAccount(address, account)
	}

	spec.setPrecompile(1, &alsofGenesisSpecBuiltin{Name: "ecrecover",
		Linear: &alsofGenesisSpecLinearPricing{Base: 3000}})
	spec.setPrecompile(2, &alsofGenesisSpecBuiltin{Name: "sha256",
		Linear: &alsofGenesisSpecLinearPricing{Base: 60, Word: 12}})
	spec.setPrecompile(3, &alsofGenesisSpecBuiltin{Name: "ripemd160",
		Linear: &alsofGenesisSpecLinearPricing{Base: 600, Word: 120}})
	spec.setPrecompile(4, &alsofGenesisSpecBuiltin{Name: "identity",
		Linear: &alsofGenesisSpecLinearPricing{Base: 15, Word: 3}})
	if genesis.Config.ByzantiumBlock != nil {
		spec.setPrecompile(5, &alsofGenesisSpecBuiltin{Name: "modexp",
			StartingBlock: (hexutil.Uint64)(genesis.Config.ByzantiumBlock.Uint64())})
		spec.setPrecompile(6, &alsofGenesisSpecBuiltin{Name: "alt_bn128_G1_add",
			StartingBlock: (hexutil.Uint64)(genesis.Config.ByzantiumBlock.Uint64()),
			Linear:        &alsofGenesisSpecLinearPricing{Base: 500}})
		spec.setPrecompile(7, &alsofGenesisSpecBuiltin{Name: "alt_bn128_G1_mul",
			StartingBlock: (hexutil.Uint64)(genesis.Config.ByzantiumBlock.Uint64()),
			Linear:        &alsofGenesisSpecLinearPricing{Base: 40000}})
		spec.setPrecompile(8, &alsofGenesisSpecBuiltin{Name: "alt_bn128_pairing_product",
			StartingBlock: (hexutil.Uint64)(genesis.Config.ByzantiumBlock.Uint64())})
	}
	return spec, nil
}

func (spec *alsofGenesisSpec) setPrecompile(address byte, data *alsofGenesisSpecBuiltin) {
	if spec.Accounts == nil {
		spec.Accounts = make(map[common.UnprefixedAddress]*alsofGenesisSpecAccount)
	}
	addr := common.UnprefixedAddress(common.BytesToAddress([]byte{address}))
	if _, exist := spec.Accounts[addr]; !exist {
		spec.Accounts[addr] = &alsofGenesisSpecAccount{}
	}
	spec.Accounts[addr].Precompiled = data
}

func (spec *alsofGenesisSpec) setAccount(address common.Address, account core.GenesisAccount) {
	if spec.Accounts == nil {
		spec.Accounts = make(map[common.UnprefixedAddress]*alsofGenesisSpecAccount)
	}

	a, exist := spec.Accounts[common.UnprefixedAddress(address)]
	if !exist {
		a = &alsofGenesisSpecAccount{}
		spec.Accounts[common.UnprefixedAddress(address)] = a
	}
	a.Balance = (*math2.HexOrDecimal256)(account.Balance)
	a.Nonce = account.Nonce

}

func (spec *alsofGenesisSpec) setByzantium(num *big.Int) {
	spec.Params.ByzantiumForkBlock = hexutil.Uint64(num.Uint64())
}

func (spec *alsofGenesisSpec) setConstantinople(num *big.Int) {
	spec.Params.ConstantinopleForkBlock = hexutil.Uint64(num.Uint64())
}

// susyChainSpec is the chain specification format used by Susy.
type susyChainSpec struct {
	Name    string `json:"name"`
	Datadir string `json:"dataDir"`
	Engine  struct {
		Sofash struct {
			Params struct {
				MinimumDifficulty      *hexutil.Big      `json:"minimumDifficulty"`
				DifficultyBoundDivisor *hexutil.Big      `json:"difficultyBoundDivisor"`
				DurationLimit          *hexutil.Big      `json:"durationLimit"`
				BlockReward            map[string]string `json:"blockReward"`
				DifficultyBombDelays   map[string]string `json:"difficultyBombDelays"`
				HomesteadTransition    hexutil.Uint64    `json:"homesteadTransition"`
				SIP100bTransition      hexutil.Uint64    `json:"sip100bTransition"`
			} `json:"params"`
		} `json:"Sofash"`
	} `json:"engine"`

	Params struct {
		AccountStartNonce        hexutil.Uint64       `json:"accountStartNonce"`
		MaximumExtraDataSize     hexutil.Uint64       `json:"maximumExtraDataSize"`
		MinGasLimit              hexutil.Uint64       `json:"minGasLimit"`
		GasLimitBoundDivisor     math2.HexOrDecimal64 `json:"gasLimitBoundDivisor"`
		NetworkID                hexutil.Uint64       `json:"networkID"`
		ChainID                  hexutil.Uint64       `json:"chainID"`
		MaxCodeSize              hexutil.Uint64       `json:"maxCodeSize"`
		MaxCodeSizeTransition    hexutil.Uint64       `json:"maxCodeSizeTransition"`
		SIP98Transition          hexutil.Uint64       `json:"sip98Transition"`
		SIP150Transition         hexutil.Uint64       `json:"sip150Transition"`
		SIP160Transition         hexutil.Uint64       `json:"sip160Transition"`
		SIP161abcTransition      hexutil.Uint64       `json:"sip161abcTransition"`
		SIP161dTransition        hexutil.Uint64       `json:"sip161dTransition"`
		SIP155Transition         hexutil.Uint64       `json:"sip155Transition"`
		SIP140Transition         hexutil.Uint64       `json:"sip140Transition"`
		SIP211Transition         hexutil.Uint64       `json:"sip211Transition"`
		SIP214Transition         hexutil.Uint64       `json:"sip214Transition"`
		SIP658Transition         hexutil.Uint64       `json:"sip658Transition"`
		SIP145Transition         hexutil.Uint64       `json:"sip145Transition"`
		SIP1014Transition        hexutil.Uint64       `json:"sip1014Transition"`
		SIP1052Transition        hexutil.Uint64       `json:"sip1052Transition"`
		SIP1283Transition        hexutil.Uint64       `json:"sip1283Transition"`
		SIP1283DisableTransition hexutil.Uint64       `json:"sip1283DisableTransition"`
	} `json:"params"`

	Genesis struct {
		Seal struct {
			Sophon struct {
				Nonce   hexutil.Bytes `json:"nonce"`
				MixHash hexutil.Bytes `json:"mixHash"`
			} `json:"sophon"`
		} `json:"seal"`

		Difficulty *hexutil.Big   `json:"difficulty"`
		Author     common.Address `json:"author"`
		Timestamp  hexutil.Uint64 `json:"timestamp"`
		ParentHash common.Hash    `json:"parentHash"`
		ExtraData  hexutil.Bytes  `json:"extraData"`
		GasLimit   hexutil.Uint64 `json:"gasLimit"`
	} `json:"genesis"`

	Nodes    []string                                             `json:"nodes"`
	Accounts map[common.UnprefixedAddress]*susyChainSpecAccount `json:"accounts"`
}

// susyChainSpecAccount is the prefunded genesis account and/or precompiled
// contract definition.
type susyChainSpecAccount struct {
	Balance math2.HexOrDecimal256   `json:"balance"`
	Nonce   math2.HexOrDecimal64    `json:"nonce,omitempty"`
	Builtin *susyChainSpecBuiltin `json:"builtin,omitempty"`
}

// susyChainSpecBuiltin is the precompiled contract definition.
type susyChainSpecBuiltin struct {
	Name       string                  `json:"name,omitempty"`
	ActivateAt math2.HexOrDecimal64    `json:"activate_at,omitempty"`
	Pricing    *susyChainSpecPricing `json:"pricing,omitempty"`
}

// susyChainSpecPricing represents the different pricing models that builtin
// contracts might advertise using.
type susyChainSpecPricing struct {
	Linear       *susyChainSpecLinearPricing       `json:"linear,omitempty"`
	ModExp       *susyChainSpecModExpPricing       `json:"modexp,omitempty"`
	AltBnPairing *susyChainSpecAltBnPairingPricing `json:"alt_bn128_pairing,omitempty"`
}

type susyChainSpecLinearPricing struct {
	Base uint64 `json:"base"`
	Word uint64 `json:"word"`
}

type susyChainSpecModExpPricing struct {
	Divisor uint64 `json:"divisor"`
}

type susyChainSpecAltBnPairingPricing struct {
	Base uint64 `json:"base"`
	Pair uint64 `json:"pair"`
}

// newSusyChainSpec converts a susy-graviton genesis block into a Susy specific
// chain specification format.
func newSusyChainSpec(network string, genesis *core.Genesis, bootnodes []string) (*susyChainSpec, error) {
	// Only sofash is currently supported between susy-graviton and Susy
	if genesis.Config.Sofash == nil {
		return nil, errors.New("unsupported consensus engine")
	}
	// Reconstruct the chain spec in Susy's format
	spec := &susyChainSpec{
		Name:    network,
		Nodes:   bootnodes,
		Datadir: strings.ToLower(network),
	}
	spec.Engine.Sofash.Params.BlockReward = make(map[string]string)
	spec.Engine.Sofash.Params.DifficultyBombDelays = make(map[string]string)
	// Frontier
	spec.Engine.Sofash.Params.MinimumDifficulty = (*hexutil.Big)(params.MinimumDifficulty)
	spec.Engine.Sofash.Params.DifficultyBoundDivisor = (*hexutil.Big)(params.DifficultyBoundDivisor)
	spec.Engine.Sofash.Params.DurationLimit = (*hexutil.Big)(params.DurationLimit)
	spec.Engine.Sofash.Params.BlockReward["0x0"] = hexutil.EncodeBig(sofash.FrontierBlockReward)

	// Homestead
	spec.Engine.Sofash.Params.HomesteadTransition = hexutil.Uint64(genesis.Config.HomesteadBlock.Uint64())

	// Tangerine Whistle : 150
	// https://github.com/susy-go/SIPs/blob/master/SIPS/sip-608.md
	spec.Params.SIP150Transition = hexutil.Uint64(genesis.Config.SIP150Block.Uint64())

	// Spurious Dragon: 155, 160, 161, 170
	// https://github.com/susy-go/SIPs/blob/master/SIPS/sip-607.md
	spec.Params.SIP155Transition = hexutil.Uint64(genesis.Config.SIP155Block.Uint64())
	spec.Params.SIP160Transition = hexutil.Uint64(genesis.Config.SIP155Block.Uint64())
	spec.Params.SIP161abcTransition = hexutil.Uint64(genesis.Config.SIP158Block.Uint64())
	spec.Params.SIP161dTransition = hexutil.Uint64(genesis.Config.SIP158Block.Uint64())

	// Byzantium
	if num := genesis.Config.ByzantiumBlock; num != nil {
		spec.setByzantium(num)
	}
	// Constantinople
	if num := genesis.Config.ConstantinopleBlock; num != nil {
		spec.setConstantinople(num)
	}
	// ConstantinopleFix (remove sip-1283)
	if num := genesis.Config.PetersburgBlock; num != nil {
		spec.setConstantinopleFix(num)
	}

	spec.Params.MaximumExtraDataSize = (hexutil.Uint64)(params.MaximumExtraDataSize)
	spec.Params.MinGasLimit = (hexutil.Uint64)(params.MinGasLimit)
	spec.Params.GasLimitBoundDivisor = (math2.HexOrDecimal64)(params.GasLimitBoundDivisor)
	spec.Params.NetworkID = (hexutil.Uint64)(genesis.Config.ChainID.Uint64())
	spec.Params.ChainID = (hexutil.Uint64)(genesis.Config.ChainID.Uint64())
	spec.Params.MaxCodeSize = params.MaxCodeSize
	// graviton has it set from zero
	spec.Params.MaxCodeSizeTransition = 0

	// Disable this one
	spec.Params.SIP98Transition = math.MaxInt64

	spec.Genesis.Seal.Sophon.Nonce = (hexutil.Bytes)(make([]byte, 8))
	binary.LittleEndian.PutUint64(spec.Genesis.Seal.Sophon.Nonce[:], genesis.Nonce)

	spec.Genesis.Seal.Sophon.MixHash = (hexutil.Bytes)(genesis.Mixhash[:])
	spec.Genesis.Difficulty = (*hexutil.Big)(genesis.Difficulty)
	spec.Genesis.Author = genesis.Coinbase
	spec.Genesis.Timestamp = (hexutil.Uint64)(genesis.Timestamp)
	spec.Genesis.ParentHash = genesis.ParentHash
	spec.Genesis.ExtraData = (hexutil.Bytes)(genesis.ExtraData)
	spec.Genesis.GasLimit = (hexutil.Uint64)(genesis.GasLimit)

	spec.Accounts = make(map[common.UnprefixedAddress]*susyChainSpecAccount)
	for address, account := range genesis.Alloc {
		bal := math2.HexOrDecimal256(*account.Balance)

		spec.Accounts[common.UnprefixedAddress(address)] = &susyChainSpecAccount{
			Balance: bal,
			Nonce:   math2.HexOrDecimal64(account.Nonce),
		}
	}
	spec.setPrecompile(1, &susyChainSpecBuiltin{Name: "ecrecover",
		Pricing: &susyChainSpecPricing{Linear: &susyChainSpecLinearPricing{Base: 3000}}})

	spec.setPrecompile(2, &susyChainSpecBuiltin{
		Name: "sha256", Pricing: &susyChainSpecPricing{Linear: &susyChainSpecLinearPricing{Base: 60, Word: 12}},
	})
	spec.setPrecompile(3, &susyChainSpecBuiltin{
		Name: "ripemd160", Pricing: &susyChainSpecPricing{Linear: &susyChainSpecLinearPricing{Base: 600, Word: 120}},
	})
	spec.setPrecompile(4, &susyChainSpecBuiltin{
		Name: "identity", Pricing: &susyChainSpecPricing{Linear: &susyChainSpecLinearPricing{Base: 15, Word: 3}},
	})
	if genesis.Config.ByzantiumBlock != nil {
		blnum := math2.HexOrDecimal64(genesis.Config.ByzantiumBlock.Uint64())
		spec.setPrecompile(5, &susyChainSpecBuiltin{
			Name: "modexp", ActivateAt: blnum, Pricing: &susyChainSpecPricing{ModExp: &susyChainSpecModExpPricing{Divisor: 20}},
		})
		spec.setPrecompile(6, &susyChainSpecBuiltin{
			Name: "alt_bn128_add", ActivateAt: blnum, Pricing: &susyChainSpecPricing{Linear: &susyChainSpecLinearPricing{Base: 500}},
		})
		spec.setPrecompile(7, &susyChainSpecBuiltin{
			Name: "alt_bn128_mul", ActivateAt: blnum, Pricing: &susyChainSpecPricing{Linear: &susyChainSpecLinearPricing{Base: 40000}},
		})
		spec.setPrecompile(8, &susyChainSpecBuiltin{
			Name: "alt_bn128_pairing", ActivateAt: blnum, Pricing: &susyChainSpecPricing{AltBnPairing: &susyChainSpecAltBnPairingPricing{Base: 100000, Pair: 80000}},
		})
	}
	return spec, nil
}

func (spec *susyChainSpec) setPrecompile(address byte, data *susyChainSpecBuiltin) {
	if spec.Accounts == nil {
		spec.Accounts = make(map[common.UnprefixedAddress]*susyChainSpecAccount)
	}
	a := common.UnprefixedAddress(common.BytesToAddress([]byte{address}))
	if _, exist := spec.Accounts[a]; !exist {
		spec.Accounts[a] = &susyChainSpecAccount{}
	}
	spec.Accounts[a].Builtin = data
}

func (spec *susyChainSpec) setByzantium(num *big.Int) {
	spec.Engine.Sofash.Params.BlockReward[hexutil.EncodeBig(num)] = hexutil.EncodeBig(sofash.ByzantiumBlockReward)
	spec.Engine.Sofash.Params.DifficultyBombDelays[hexutil.EncodeBig(num)] = hexutil.EncodeUint64(3000000)
	n := hexutil.Uint64(num.Uint64())
	spec.Engine.Sofash.Params.SIP100bTransition = n
	spec.Params.SIP140Transition = n
	spec.Params.SIP211Transition = n
	spec.Params.SIP214Transition = n
	spec.Params.SIP658Transition = n
}

func (spec *susyChainSpec) setConstantinople(num *big.Int) {
	spec.Engine.Sofash.Params.BlockReward[hexutil.EncodeBig(num)] = hexutil.EncodeBig(sofash.ConstantinopleBlockReward)
	spec.Engine.Sofash.Params.DifficultyBombDelays[hexutil.EncodeBig(num)] = hexutil.EncodeUint64(2000000)
	n := hexutil.Uint64(num.Uint64())
	spec.Params.SIP145Transition = n
	spec.Params.SIP1014Transition = n
	spec.Params.SIP1052Transition = n
	spec.Params.SIP1283Transition = n
}

func (spec *susyChainSpec) setConstantinopleFix(num *big.Int) {
	spec.Params.SIP1283DisableTransition = hexutil.Uint64(num.Uint64())
}

// pySophonGenesisSpec represents the genesis specification format used by the
// Python Sophon implementation.
type pySophonGenesisSpec struct {
	Nonce      hexutil.Bytes     `json:"nonce"`
	Timestamp  hexutil.Uint64    `json:"timestamp"`
	ExtraData  hexutil.Bytes     `json:"extraData"`
	GasLimit   hexutil.Uint64    `json:"gasLimit"`
	Difficulty *hexutil.Big      `json:"difficulty"`
	Mixhash    common.Hash       `json:"mixhash"`
	Coinbase   common.Address    `json:"coinbase"`
	Alloc      core.GenesisAlloc `json:"alloc"`
	ParentHash common.Hash       `json:"parentHash"`
}

// newPySophonGenesisSpec converts a susy-graviton genesis block into a Susy specific
// chain specification format.
func newPySophonGenesisSpec(network string, genesis *core.Genesis) (*pySophonGenesisSpec, error) {
	// Only sofash is currently supported between susy-graviton and pysophon
	if genesis.Config.Sofash == nil {
		return nil, errors.New("unsupported consensus engine")
	}
	spec := &pySophonGenesisSpec{
		Timestamp:  (hexutil.Uint64)(genesis.Timestamp),
		ExtraData:  genesis.ExtraData,
		GasLimit:   (hexutil.Uint64)(genesis.GasLimit),
		Difficulty: (*hexutil.Big)(genesis.Difficulty),
		Mixhash:    genesis.Mixhash,
		Coinbase:   genesis.Coinbase,
		Alloc:      genesis.Alloc,
		ParentHash: genesis.ParentHash,
	}
	spec.Nonce = (hexutil.Bytes)(make([]byte, 8))
	binary.LittleEndian.PutUint64(spec.Nonce[:], genesis.Nonce)

	return spec, nil
}
