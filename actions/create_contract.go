// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package actions

import (
	"context"
	"encoding/binary"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/hypersdk/chain"
	"github.com/ava-labs/hypersdk/codec"
	"github.com/ava-labs/hypersdk/consts"
	mconsts "github.com/ava-labs/hypersdk/examples/morpheusvm/consts"
	"github.com/ava-labs/hypersdk/examples/morpheusvm/storage"
	"github.com/ava-labs/hypersdk/state"
	"github.com/ava-labs/hypersdk/utils"
)

var _ chain.Action = (*CreateContract)(nil)

type CreateContract struct {
	Bytecode      []byte
	InitialState  []byte
	Discriminator uint16
}

func (*CreateContract) GetTypeID() uint8 {
	return mconsts.CreateContractID
}

func (t *CreateContract) StateKeys(actor codec.Address, _ ids.ID) state.Keys {
	contractAddress := storage.GenerateContractAddress(actor, t.Discriminator)

	return state.Keys{
		string(storage.ContractStateKey(contractAddress)):    state.All,
		string(storage.ContractBytecodeKey(contractAddress)): state.All,
	}
}

func (*CreateContract) StateKeysMaxChunks() []uint16 {
	return []uint16{storage.ContractStateChunks, storage.ContractBytecodeChunks}
}

func (*CreateContract) OutputsWarpMessage() bool {
	return false
}

func (t *CreateContract) Execute(
	ctx context.Context,
	_ chain.Rules,
	mu state.Mutable,
	_ int64,
	actor codec.Address,
	_ ids.ID,
	_ bool,
) (bool, uint64, []byte, *warp.UnsignedMessage, error) {
	addr, err := storage.CreateContract(ctx, mu, actor, t.Bytecode, t.InitialState, t.Discriminator)
	if err != nil {
		return false, 1, utils.ErrBytes(err), nil, nil
	}

	addrString := codec.MustAddressBech32(mconsts.HRP,addr)

	//success, computeUnits, output, warpMsg, err
	return true, 1, []byte(addrString), nil, nil
}

func (*CreateContract) MaxComputeUnits(chain.Rules) uint64 {
	return 1
}

func (cc *CreateContract) Size() int {
	return len(cc.Bytecode) + len(cc.InitialState) + consts.Uint8Len
}

func (t *CreateContract) Marshal(p *codec.Packer) {
	p.PackBytes(t.Bytecode)
	p.PackBytes(t.InitialState)

	discriminatorBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(discriminatorBytes, t.Discriminator)
	p.PackBytes(discriminatorBytes)
}

func UnmarshalCreateContract(p *codec.Packer, _ *warp.Message) (chain.Action, error) {
	var action CreateContract
	p.UnpackBytes(-1, false, &action.Bytecode)
	p.UnpackBytes(-1, false, &action.InitialState)

	var discriminatorBytes []byte = make([]byte, 2)
	p.UnpackBytes(2, false, &discriminatorBytes)
	action.Discriminator = binary.BigEndian.Uint16(discriminatorBytes)

	return &action, nil
}

func (*CreateContract) ValidRange(chain.Rules) (int64, int64) {
	// Returning -1, -1 means that the action is always valid.
	return -1, -1
}
