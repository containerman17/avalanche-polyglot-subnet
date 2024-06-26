// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rpc

import (
	"context"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/trace"
	"github.com/ava-labs/hypersdk/codec"
	"github.com/ava-labs/hypersdk/fees"
	"github.com/containerman17/avalanche-polyglot-subnet/genesis"
)

type Controller interface {
	Genesis() *genesis.Genesis
	Tracer() trace.Tracer
	GetTransaction(context.Context, ids.ID) (bool, int64, bool, fees.Dimensions, uint64, error)
	GetBalanceFromState(context.Context, codec.Address) (uint64, error)
	GetContractBytecodeFromState(context.Context, codec.Address) ([]byte, error)
	GetContractStateFromState(context.Context, codec.Address) ([]byte, error)
}
