import { Uint8ArrayToBase64, base64ToUint8Array, readStdin, writeStdOut } from "./javy_io";
import {
    type deserialize as BorshDeserialize,
    type serialize as BorshSerialize,
    type validate as BorshValidate,
} from "@dao-xyz/borsh";


export { Uint8ArrayToBase64, base64ToUint8Array };

export abstract class FunctionCallParams {
}

export interface ContractIFace {
    execute(callParams: FunctionCallParams, actor: Uint8Array): any
}

export function executeContract<T extends ContractIFace>(
    ContractClass: new () => T,
    serialize: typeof BorshSerialize,
    deserialize: typeof BorshDeserialize,
    validate: typeof BorshValidate
) {
    try {
        type incomingMessage = {
            currentState: string,
            payload: string,
            actor: string,
        }

        const stdinStr = readStdin()

        const input = JSON.parse(stdinStr) as incomingMessage
        const contractStateBytes: Uint8Array = base64ToUint8Array(input.currentState || "")
        const payloadBytes: Uint8Array = base64ToUint8Array(input.payload || "")

        let contractInstance: T;
        if (contractStateBytes.length === 0) {
            console.log(`Initializing an empty contract`)

            validate(ContractClass)
            contractInstance = new ContractClass();
        } else {
            console.log(`Deserializing contract state`)
            contractInstance = deserialize(contractStateBytes, ContractClass);
        }

        const callParamsDecoded = deserialize(payloadBytes, FunctionCallParams)

        const result = contractInstance.execute(callParamsDecoded, base64ToUint8Array(input.actor || ""))

        const endState = serialize(contractInstance)

        writeStdOut(JSON.stringify({
            success: true,
            endState: Uint8ArrayToBase64(endState),
            result: Uint8ArrayToBase64(result ? serialize(result) : new Uint8Array(0)),
        }))
    } catch (e) {
        writeStdOut(JSON.stringify({
            success: false,
            error: String(e) + "\n" + (e as Error).stack,
        }))
    }
}

