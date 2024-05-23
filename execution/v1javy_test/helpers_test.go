package v1javy_test

//go:generate bash -c "cd assets && npm ci"
//go:generate npx ../v1javy/js_sdk/ assets/counters.ts

import (
	"log"

	_ "embed"

	"github.com/ava-labs/hypersdk/codec"
	"github.com/near/borsh-go"
)

//go:embed assets/counters.wasm
var testWasmBytes []byte

//MySuperCalculator.loadCPU

type loadCPUCallParams struct {
	n uint16
}

const LOAD_CPU_VARIANT = 0x3

func generateLoadCPUPayload(n uint16) []byte {
	bytes, err := borsh.Serialize(loadCPUCallParams{
		n: n,
	})
	if err != nil {
		log.Fatal(err)
	}

	return append([]byte{LOAD_CPU_VARIANT}, bytes...)
}

//MySuperCalculator.increment

const INCREMENT_VARIANT = 0x0

type incrementCallParams struct {
	Amount uint64
}

func generateIncrementPayload(amount uint64) []byte {
	bytes, err := borsh.Serialize(incrementCallParams{
		Amount: amount,
	})
	if err != nil {
		log.Fatal(err)
	}

	return append([]byte{INCREMENT_VARIANT}, bytes...)
}

//MySuperCalculator.getCounter

const GET_COUNTER_VARIANT = 0x2

type getCounterCallParams struct {
	User [33]byte
}

func generateGetCounterPayload(user []byte) []byte {
	fixedUser := [33]byte{}
	copy(fixedUser[:], user)

	bytes, err := borsh.Serialize(getCounterCallParams{
		User: fixedUser,
	})
	if err != nil {
		log.Fatal(err)
	}

	return append([]byte{GET_COUNTER_VARIANT}, bytes...)
}

type getCounterResult struct {
	Counter uint64
}

func decodeGetCounterResult(payload []byte) uint64 {
	result := &getCounterResult{}
	err := borsh.Deserialize(result, payload)
	if err != nil {
		log.Fatal(err)
	}

	return result.Counter
}

func createActorAddress(actorNumber uint) []byte {
	actor := codec.Address{byte(actorNumber), byte(actorNumber), byte(actorNumber)}
	actorBytes := make([]byte, len(actor))
	copy(actorBytes, actor[:])
	return actorBytes
}
