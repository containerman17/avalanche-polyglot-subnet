//WARNING: if changed, you have to call go generate ./...
import { deserialize, field, fixedArray, serialize, validate, variant, vec } from "@dao-xyz/borsh";
import { executeContract, ContractIFace, FunctionCallParams } from "../../v1javy/js_sdk/src/index";

class AddressToUint64Map {
    //an array of 33-byte arrays
    @field({ type: vec(fixedArray('u8', 33)) })
    public keys: Uint8Array[] = [];

    @field({ type: vec('u64') })
    public values: bigint[] = [];

    getValue(key: Uint8Array): bigint {
        const index = this.keys.findIndex((k) => k.every((byte, i) => byte === key[i]));
        if (index === -1) {
            return BigInt(0);
        }
        return this.values[index];
    }

    setValue(key: Uint8Array, value: bigint) {
        const index = this.keys.findIndex((k) => k.every((byte, i) => byte === key[i]));
        if (index === -1) {
            this.keys.push(key);
            this.values.push(value);
        } else {
            this.values[index] = value;
        }
    }
}

class MySuperCalculator implements ContractIFace {
    @field({ type: AddressToUint64Map })
    userCounters: AddressToUint64Map = new AddressToUint64Map();

    execute(callParams: FunctionCallParams, actor: Uint8Array) {
        switch (callParams.constructor) {
            case IncrementCallParams:
                return this.increment(callParams as IncrementCallParams, actor);
            case DecrementCallParams:
                return this.decrement(callParams as DecrementCallParams, actor);
            case GetCounterCallParams:
                return this.getCounter(callParams as GetCounterCallParams, actor);
            case LoadCPUCallParams:
                return this.loadCPU(callParams as LoadCPUCallParams, actor);
            default:
                throw new Error("Unknown callParams type");
        }
    }

    increment({ amount }: IncrementCallParams, actor: Uint8Array) {
        this.userCounters.setValue(actor, this.userCounters.getValue(actor) + amount);
    }

    decrement({ amount }: DecrementCallParams, actor: Uint8Array) {
        this.userCounters.setValue(actor, this.userCounters.getValue(actor) - amount);
    }

    loadCPU({ n }: LoadCPUCallParams, _: Uint8Array) {
        let summ = 0;
        for (let i = 0; i < n; i++) {
            for (let j = 0; j < n; j++) {
                summ += i * j;
            }
        }
    }

    getCounter({ user }: GetCounterCallParams, _: Uint8Array): GetCounterResult {
        const result = new GetCounterResult();
        result.counter = this.userCounters.getValue(user);
        return result;
    }
}

@variant(0)
class IncrementCallParams extends FunctionCallParams {
    @field({ type: "u64" })
    public amount: bigint = BigInt(0)
}

@variant(1)
class DecrementCallParams extends FunctionCallParams {
    @field({ type: "u64" })
    public amount: bigint = BigInt(0)
}

@variant(2)
class GetCounterCallParams extends FunctionCallParams {
    @field({ type: fixedArray('u8', 33) })
    public user: Uint8Array = new Uint8Array(33);
}

@variant(3)
class LoadCPUCallParams extends FunctionCallParams {
    @field({ type: "u16" })
    public n: number = 0;
}

//return types do not need @variant
class GetCounterResult {
    @field({ type: "u64" })
    public counter: bigint = BigInt(0);
}

executeContract(MySuperCalculator, serialize, deserialize, validate);
