package v1javy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bytecodealliance/wasmtime-go/v19"
)

type JavyExecParams struct {
	MaxFuel      uint64        `json:"-"`
	MaxTime      time.Duration `json:"-"`
	MaxMemory    int64         `json:"-"`
	Bytecode     *[]byte       `json:"-"`
	CurrentState []byte        `json:"currentState"`
	Payload      []byte        `json:"payload"`
	Actor        []byte        `json:"actor"`
}

type JavyExecResult struct {
	FuelConsumed uint64
	TimeTaken    time.Duration
	UpdatedState *[]byte //nil if no update
	StdErr       []byte
	Result       []byte
}

type JavyExec struct {
	executeMutexes map[uint64]*sync.Mutex
	storesCache    map[uint64]*struct {
		store    *wasmtime.Store
		mainFunc *wasmtime.Func
	}
}

func NewJavyExec() *JavyExec {
	return &JavyExec{
		executeMutexes: map[uint64]*sync.Mutex{},
		storesCache: map[uint64]*struct {
			store    *wasmtime.Store
			mainFunc *wasmtime.Func
		}{},
	}
}

type stdoutResultJson struct {
	Result   []byte `json:"result"`
	Success  bool   `json:"success"`
	EndState []byte `json:"endState"`
	Error    string `json:"error"`
}

func (exec *JavyExec) Execute(params JavyExecParams) (*JavyExecResult, error) {
	store, mainFunc, err := exec.createStore(*params.Bytecode)
	if err != nil {
		return nil, err
	}

	store.Limiter(params.MaxMemory, -1, -1, -1, -1)
	defer store.Engine.Close()
	defer store.Close()

	callDataJson, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshalling call data: %v", err)
	}

	//FIXME: use syscall.Mkfifo instead of torturing disk with temp files
	stdoutFile, err := os.CreateTemp("", "stdout*")
	if err != nil {
		return nil, fmt.Errorf("creating stdout file: %v", err)
	}
	defer os.Remove(stdoutFile.Name())

	stderrFile, err := os.CreateTemp("", "stderr*")
	if err != nil {
		return nil, fmt.Errorf("creating stderr file: %v", err)
	}
	defer os.Remove(stderrFile.Name())

	stdinFile, err := os.CreateTemp("", "stdin*")
	if err != nil {
		return nil, fmt.Errorf("creating stdin file: %v", err)
	}
	defer os.Remove(stdinFile.Name())

	_, err = stdinFile.Write(callDataJson)
	if err != nil {
		return nil, fmt.Errorf("writing to stdin file: %v", err)
	}

	err = store.SetFuel(params.MaxFuel)
	if err != nil {
		return nil, fmt.Errorf("setting fuel: %v", err)
	}

	wasiConfig := wasmtime.NewWasiConfig()
	wasiConfig.SetStdoutFile(stdoutFile.Name())
	wasiConfig.SetStderrFile(stderrFile.Name())
	wasiConfig.SetStdinFile(stdinFile.Name())
	store.SetWasi(wasiConfig)

	startTime := time.Now()

	finished := false

	timeoutErrCh := make(chan error, 1)
	go func() {
		time.Sleep(params.MaxTime)
		if !finished {
			fmt.Printf("Execution timed out\n")
			store.Engine.IncrementEpoch()
			timeoutErrCh <- fmt.Errorf("execution timed out")
		}
	}()

	_, err = mainFunc.Call(store)
	finished = true
	if err != nil {
		return nil, fmt.Errorf("calling user code main function: %v", err)
	}
	execTime := time.Since(startTime)

	select {
	case err := <-timeoutErrCh:
		if err != nil {
			return nil, err
		}
	default:
	}

	fuelAfter, err := store.GetFuel()
	if err != nil {
		return nil, fmt.Errorf("getting fuel after execution: %v", err)
	}
	consumedFuel := params.MaxFuel - fuelAfter

	stdoutBytes, err := os.ReadFile(stdoutFile.Name())
	if err != nil {
		return nil, fmt.Errorf("reading stdout file: %v", err)
	}

	var stdoutResult stdoutResultJson
	err = json.Unmarshal(stdoutBytes, &stdoutResult)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling stdout: %v", err)
	}
	if !stdoutResult.Success {
		return nil, fmt.Errorf("execution failed: %s", string(stdoutBytes))
	}

	stderrBytes, err := os.ReadFile(stderrFile.Name())
	if err != nil {
		return nil, fmt.Errorf("reading stderr file: %v", err)
	}

	var updatedState *[]byte = nil
	if !bytes.Equal(stdoutResult.EndState, params.CurrentState) {
		updatedState = &stdoutResult.EndState
	}

	return &JavyExecResult{
		FuelConsumed: consumedFuel,
		TimeTaken:    execTime,
		UpdatedState: updatedState,
		StdErr:       stderrBytes,
		Result:       stdoutResult.Result,
	}, nil
}
