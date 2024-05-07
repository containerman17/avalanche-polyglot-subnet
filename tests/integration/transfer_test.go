package integration_test

import (
	"context"
	"testing"

	"github.com/ava-labs/hypersdk/examples/morpheusvm/actions"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	prep := prepare(t)
	//check sender's initial balance is 10kk
	senderInitBalance, err := prep.instance.lcli.Balance(context.Background(), prep.addrStr)
	require.NoError(t, err)
	require.Equal(t, uint64(10_000_000), senderInitBalance)

	//check receiver's initial balance is 0
	receiverInitBalance, err := prep.instance.lcli.Balance(context.Background(), prep.addrStr2)
	require.NoError(t, err)
	require.Equal(t, uint64(0), receiverInitBalance)

	//issue Transfer transaction
	parser, err := prep.instance.lcli.Parser(context.Background())
	require.NoError(t, err)
	submit, _, _, err := prep.instance.cli.GenerateTransaction(
		context.Background(),
		parser,
		nil,
		&actions.Transfer{
			Value: 500_000,    // transfer amount
			To:    prep.addr2, // receiver address
		},
		prep.factory,
	)
	require.NoError(t, err)
	require.NoError(t, submit(context.Background()))

	results := prep.expectBlk(t, prep.instance)(false)
	require.Len(t, results, 1)
	require.True(t, results[0].Success)

	const EXPECTED_GAS_FEE = 261 //FIXME:

	//check the final balances
	senderFinalBalance, err := prep.instance.lcli.Balance(context.Background(), prep.addrStr)
	require.NoError(t, err)
	expectedSenderFinalBalance := uint64(10_000_000 - 500_000 - EXPECTED_GAS_FEE)
	require.Equal(t, expectedSenderFinalBalance, senderFinalBalance, "Expected sender final balance: %d, got: %d", expectedSenderFinalBalance, senderFinalBalance)

	receiverFinalBalance, err := prep.instance.lcli.Balance(context.Background(), prep.addrStr2)
	require.NoError(t, err)
	expectedReceiverFinalBalance := uint64(500_000)
	require.Equal(t, expectedReceiverFinalBalance, receiverFinalBalance, "Expected receiver final balance: %d, got: %d", expectedReceiverFinalBalance, receiverFinalBalance)
}
