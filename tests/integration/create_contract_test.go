package integration_test

import (
	"context"
	"testing"

	"github.com/containerman17/avalanche-polyglot-subnet/actions"
	"github.com/stretchr/testify/require"
)

func TestCreateContract(t *testing.T) {
	dummyBytecode := []byte{0x01, 0x02, 0x03}
	dummyState := []byte{0x04, 0x05, 0x06}
	discriminator := 123

	//send CreateContract tx

	prep := prepare(t)

	parser, err := prep.instance.lcli.Parser(context.Background())
	require.NoError(t, err)
	submit, _, _, err := prep.instance.cli.GenerateTransaction(
		context.Background(),
		parser,
		nil,
		&actions.CreateContract{
			Bytecode:      dummyBytecode,
			InitialState:  dummyState,
			Discriminator: uint16(discriminator),
		},
		prep.factory,
	)
	require.NoError(t, err)
	require.NoError(t, submit(context.Background()))

	results := prep.expectBlk(t, prep.instance)(false)
	require.Len(t, results, 1)
	require.True(t, results[0].Success)

	contractAddrString := string(results[0].Output)

	//check bytecode and state
	bytecodeFromChain, err := prep.instance.lcli.ContractBytecode(context.Background(), contractAddrString)
	require.NoError(t, err)
	require.Equal(t, []byte{0x01, 0x02, 0x03}, bytecodeFromChain)

	stateFromChain, err := prep.instance.lcli.ContractState(context.Background(), contractAddrString)
	require.NoError(t, err)
	require.Equal(t, []byte{0x04, 0x05, 0x06}, stateFromChain)

}
