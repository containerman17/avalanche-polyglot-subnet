# Polyglot Javy SDk
Polyglot Javy SDK Assists with smart contract creation and compilation for Avalanche Polyglot subnet.

## Example contract
```typescript
import { deserialize, field, serialize, variant, vec } from "@dao-xyz/borsh";
import { executeContract, ContractIFace, FunctionCallParams } from "avax-polyglot-sdk-javy";

class MySuperCalculator implements ContractIFace {
  @field({ type: 'u64' }) // uint64 requires bigint in JS
  counter: bigint = BigInt(0); // initial value matters only on the first transaction to this contract

  @field({ type: vec('string') }) // vec is a js array, check out @dao-xyz/borsh to learn more
  log: string[] = [];

  execute(callParams: FunctionCallParams) {
    // Required by ContractIFace. Routes by callParams type
    if (callParams instanceof IncrementCallParams) {
      this.increment(callParams);
    } else if (callParams instanceof DecrementCallParams) {
      this.decrement(callParams);
    }
  }

  increment({ amount, message }: IncrementCallParams) { // custom function
    this.counter += amount;
    this.log.push(message);
  }

  decrement({ amount, message }: DecrementCallParams) { // custom function
    this.counter -= amount;
    this.log.push(message);
  }
}

@variant(0)
class IncrementCallParams extends FunctionCallParams { // has to extend FunctionCallParams for serialization purposes
  @field({ type: "u64" })
  public amount: bigint = BigInt(0);

  @field({ type: "string" })
  public message: string = "";
}

@variant(1) // has to be different for every class extending FunctionCallParams
class DecrementCallParams extends FunctionCallParams {
  @field({ type: "u64" })
  public amount: bigint = BigInt(0);

  @field({ type: "string" })
  public message: string = "";
}

executeContract(MySuperCalculator, serialize, deserialize); // required
```

## Compiling to wasm bytecode
```bash
npx avax-polyglot-sdk-javy ./src/your_contract_path.ts ./dist/your_wasm_path.wasm
```

## Recommended TS-config for your project
```json
{
  "compilerOptions": {
    "target": "es2020",
    "module": "es2020",
    "moduleResolution": "node",
    "strict": true,
    "experimentalDecorators": true
  },
  "include": ["src/**/*.ts"]
}
```
`experimentalDecorators` option is required