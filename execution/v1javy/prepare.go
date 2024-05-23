package v1javy

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/bytecodealliance/wasmtime-go/v19"
)

func (exec *JavyExec) createStore(wasmBytes []byte) (*wasmtime.Store, *wasmtime.Func, error) {
	//FIXME: reusing the same store for the same smart conbtract gives 5-10x performance boost
	config := wasmtime.NewConfig()
	config.SetConsumeFuel(true)

	engine := wasmtime.NewEngineWithConfig(config)

	userCodeModule, err := wasmtime.NewModule(engine, wasmBytes) //Serialized wasm module does not add any performance benefits
	if err != nil {
		return nil, nil, fmt.Errorf("instantiating user code module: %v", err)
	}

	compiledLib, err := getCwasmBytes()
	if err != nil {
		return nil, nil, fmt.Errorf("getting javy provider compiled wasm: %v", err)
	}

	libraryModule, err := wasmtime.NewModuleDeserialize(engine, *compiledLib)
	if err != nil {
		cwasmCachePath, _ := getCwasmCachePath()
		log.Printf("Library size: %d", len(*compiledLib))
		return nil, nil, fmt.Errorf("instantiating javy library module (consider cleaning up %s): %v", cwasmCachePath, err)
	}

	store := wasmtime.NewStore(engine)

	//check out https://github.com/ava-labs/hypersdk/tree/main/x/programs/engine

	linker := wasmtime.NewLinker(engine)
	linker.DefineWasi()

	libraryInstance, err := linker.Instantiate(store, libraryModule)
	if err != nil {
		return nil, nil, fmt.Errorf("instantiating javy library instance: %v", err)
	}

	linker.DefineInstance(store, "javy_quickjs_provider_v1", libraryInstance)

	linker.AllowShadowing(true)
	userCodeInstance, err := linker.Instantiate(store, userCodeModule)
	if err != nil {
		return nil, nil, fmt.Errorf("instantiating user code instance: %v", err)
	}

	userCodeMain := userCodeInstance.GetFunc(store, "_start")

	return store, userCodeMain, nil
}
