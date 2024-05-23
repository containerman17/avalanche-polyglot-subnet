

interface JavyBuiltins {
    IO: {
        readSync(fd: number, buffer: Uint8Array): number;
        writeSync(fd: number, buffer: Uint8Array): number;
    };
}

declare global {
    const Javy: JavyBuiltins;
}

export interface SmartContract {
    serialize(): Promise<Uint8Array>;
    deSerialize(data: Uint8Array): void;
}

