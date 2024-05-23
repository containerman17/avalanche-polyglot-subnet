package v1javy_test

import (
	"testing"
	"time"

	"github.com/containerman17/avalanche-polyglot-subnet/execution/v1javy"
)

func TestReturn(t *testing.T) {
	exec := v1javy.NewJavyExec()

	actor1Bytes := createActorAddress(1)
	actor2Bytes := createActorAddress(2)
	actor3Bytes := createActorAddress(3)

	params := v1javy.JavyExecParams{
		MaxFuel:      10 * 1000 * 1000,
		MaxTime:      time.Millisecond * 20,
		MaxMemory:    1024 * 1024 * 100,
		Bytecode:     &testWasmBytes,
		CurrentState: []byte{},
		Payload:      generateIncrementPayload(131313),
	}

	//execute 2 times for actor 1
	params.Actor = actor1Bytes
	for i := 0; i < 2; i++ {
		res, err := exec.Execute(params)
		if err != nil {
			t.Fatal(err)
		}
		if res.UpdatedState != nil {
			params.CurrentState = *res.UpdatedState
		}
	}

	//execute 1 time for actor 2
	params.Actor = actor2Bytes
	res, err := exec.Execute(params)
	if err != nil {
		t.Fatal(err)
	}
	if res.UpdatedState != nil {
		params.CurrentState = *res.UpdatedState
	}

	//get result for actor 1
	params.Payload = generateGetCounterPayload(actor1Bytes)

	res, err = exec.Execute(params)
	if err != nil {
		t.Fatal(err)
	}
	if res.UpdatedState != nil {
		params.CurrentState = *res.UpdatedState
	}

	//check result for actor 1
	actor1FinalBalance := decodeGetCounterResult(res.Result)

	if actor1FinalBalance != 131313*2 {
		t.Fatalf("Expected balance %d, got %d", 131313*2, actor1FinalBalance)
	}

	//get result for actor 2
	params.Payload = generateGetCounterPayload(actor2Bytes)
	res, err = exec.Execute(params)
	if err != nil {
		t.Fatal(err)
	}
	if res.UpdatedState != nil {
		params.CurrentState = *res.UpdatedState
	}

	//check result for actor 2
	actor2FinalBalance := decodeGetCounterResult(res.Result)

	if actor2FinalBalance != 131313 {
		t.Fatalf("Expected balance %d, got %d", 131313, actor2FinalBalance)
	}

	//check zero balance for actor 3
	params.Payload = generateGetCounterPayload(actor3Bytes)
	res, err = exec.Execute(params)
	if err != nil {
		t.Fatal(err)
	}
	if res.UpdatedState != nil {
		params.CurrentState = *res.UpdatedState
	}

	actor3FinalBalance := decodeGetCounterResult(res.Result)

	if actor3FinalBalance != 0 {
		t.Fatalf("Expected balance %d, got %d", 0, actor3FinalBalance)
	}
}
