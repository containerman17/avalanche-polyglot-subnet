package v1javy_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/containerman17/avalanche-polyglot-subnet/execution/v1javy"
)

var DEFAULT_PARAMS_LIMITS = v1javy.JavyExecParams{
	MaxFuel:      10 * 1000 * 1000,
	MaxTime:      time.Millisecond * 20,
	MaxMemory:    1024 * 1024 * 100,
	Bytecode:     &testWasmBytes,
	CurrentState: []byte{},
	Payload:      generateLoadCPUPayload(60),
	Actor:        []byte{},
}

func TestMaxFuel(t *testing.T) {
	t.Parallel()

	exec := v1javy.NewJavyExec()

	params := DEFAULT_PARAMS_LIMITS

	res, err := exec.Execute(params)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("Fuel consumed: %d\n", res.FuelConsumed)

	params.MaxFuel -= res.FuelConsumed / 2

	_, err = exec.Execute(params)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestMaxTime(t *testing.T) {
	t.Parallel()

	exec := v1javy.NewJavyExec()

	params := DEFAULT_PARAMS_LIMITS

	res, err := exec.Execute(params)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("Time consumed: %v\n", res.TimeTaken)

	params.MaxTime = time.Nanosecond

	res, err = exec.Execute(params)
	if err == nil {
		fmt.Printf("Time consumed: %v\n", res.TimeTaken)
		t.FailNow()
	}
}

func TestMaxMemory(t *testing.T) {
	t.Parallel()

	exec := v1javy.NewJavyExec()

	params := DEFAULT_PARAMS_LIMITS
	params.MaxMemory = 1024 * 1024 * 5 //10MB

	_, err := exec.Execute(params)
	if err != nil {
		t.Error(err)
	}

	params.MaxMemory /= 10

	_, err = exec.Execute(params)
	if err == nil {
		t.Error("Expected error")
	}
}
