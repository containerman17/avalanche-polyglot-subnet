// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rpc

import (
	"net/http"

	"github.com/ava-labs/avalanchego/ids"

	"github.com/ava-labs/hypersdk/codec"
	"github.com/ava-labs/hypersdk/fees"
	"github.com/containerman17/avalanche-polyglot-subnet/consts"
	"github.com/containerman17/avalanche-polyglot-subnet/genesis"
)

type JSONRPCServer struct {
	c Controller
}

func NewJSONRPCServer(c Controller) *JSONRPCServer {
	return &JSONRPCServer{c}
}

type GenesisReply struct {
	Genesis *genesis.Genesis `json:"genesis"`
}

func (j *JSONRPCServer) Genesis(_ *http.Request, _ *struct{}, reply *GenesisReply) (err error) {
	reply.Genesis = j.c.Genesis()
	return nil
}

type TxArgs struct {
	TxID ids.ID `json:"txId"`
}

type TxReply struct {
	Timestamp int64           `json:"timestamp"`
	Success   bool            `json:"success"`
	Units     fees.Dimensions `json:"units"`
	Fee       uint64          `json:"fee"`
}

func (j *JSONRPCServer) Tx(req *http.Request, args *TxArgs, reply *TxReply) error {
	ctx, span := j.c.Tracer().Start(req.Context(), "Server.Tx")
	defer span.End()

	found, t, success, units, fee, err := j.c.GetTransaction(ctx, args.TxID)
	if err != nil {
		return err
	}
	if !found {
		return ErrTxNotFound
	}
	reply.Timestamp = t
	reply.Success = success
	reply.Units = units
	reply.Fee = fee
	return nil
}

type BalanceArgs struct {
	Address string `json:"address"`
}

type BalanceReply struct {
	Amount uint64 `json:"amount"`
}

func (j *JSONRPCServer) Balance(req *http.Request, args *BalanceArgs, reply *BalanceReply) error {
	ctx, span := j.c.Tracer().Start(req.Context(), "Server.Balance")
	defer span.End()

	addr, err := codec.ParseAddressBech32(consts.HRP, args.Address)
	if err != nil {
		return err
	}
	balance, err := j.c.GetBalanceFromState(ctx, addr)
	if err != nil {
		return err
	}
	reply.Amount = balance
	return err
}

type ContractBytecodeArgs struct {
	Address string `json:"address"`
}

type ContractBytecodeReply struct {
	Bytecode []byte `json:"bytecode"`
}

func (j *JSONRPCServer) ContractBytecode(req *http.Request, args *ContractBytecodeArgs, reply *ContractBytecodeReply) error {
	ctx, span := j.c.Tracer().Start(req.Context(), "Server.ContractBytecode")
	defer span.End()

	addr, err := codec.ParseAddressBech32(consts.HRP, args.Address)
	if err != nil {
		return err
	}
	bytecode, err := j.c.GetContractBytecodeFromState(ctx, addr)
	if err != nil {
		return err
	}
	reply.Bytecode = bytecode
	return err
}

type ContractStateArgs struct {
	Address string `json:"address"`
}

type ContractStateReply struct {
	State []byte `json:"state"`
}

func (j *JSONRPCServer) ContractState(req *http.Request, args *ContractStateArgs, reply *ContractStateReply) error {
	ctx, span := j.c.Tracer().Start(req.Context(), "Server.ContractState")
	defer span.End()

	addr, err := codec.ParseAddressBech32(consts.HRP, args.Address)
	if err != nil {
		return err
	}
	state, err := j.c.GetContractStateFromState(ctx, addr)
	if err != nil {
		return err
	}
	reply.State = state
	return err
}
