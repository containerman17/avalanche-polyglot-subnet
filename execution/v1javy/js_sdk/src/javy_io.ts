const enum STDIO {
    Stdin,
    Stdout,
    Stderr,
}

export function readStdin(): string {
    let buffer = new Uint8Array(1024);
    let bytesUsed = 0;
    while (true) {
        const bytesRead = Javy.IO.readSync(STDIO.Stdin, buffer.subarray(bytesUsed));
        // A negative number of bytes read indicates an error.
        if (bytesRead < 0) {
            // FIXME: Figure out the specific error that occured.
            throw Error("Error while reading from file descriptor");
        }
        // 0 bytes read means we have reached EOF.
        if (bytesRead === 0) {
            const endBuffer = buffer.subarray(0, bytesUsed + bytesRead);
            return new TextDecoder().decode(endBuffer);
        }

        bytesUsed += bytesRead;
        // If we have filled the buffer, but have not reached EOF yet,
        // double the buffers capacity and continue.
        if (bytesUsed === buffer.length) {
            const nextBuffer = new Uint8Array(buffer.length * 2);
            nextBuffer.set(buffer);
            buffer = nextBuffer;
        }
    }
}

export function writeStdOut(input: string) {
    const encoder = new TextEncoder();
    const buffer = encoder.encode(input);
    writeFileSync(STDIO.Stdout, buffer);
}

export function writeStdErr(input: string) {
    const encoder = new TextEncoder();
    const buffer = encoder.encode(input);
    writeFileSync(STDIO.Stderr, buffer);
}

function writeFileSync(fd: number, buffer: Uint8Array) {
    while (buffer.length > 0) {
        // Try to write the entire buffer.
        const bytesWritten = Javy.IO.writeSync(fd, buffer);
        // A negative number of bytes written indicates an error.
        if (bytesWritten < 0) {
            throw Error("Error while writing to file descriptor");
        }
        // 0 bytes means that the destination cannot accept additional bytes.
        if (bytesWritten === 0) {
            throw Error("Could not write all contents in buffer to file descriptor");
        }
        // Otherwise cut off the bytes from the buffer that
        // were successfully written.
        buffer = buffer.subarray(bytesWritten);
    }
}

export function base64ToUint8Array(base64String: string): Uint8Array {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';
    let str = base64String.replace(/=+$/, ''); // Remove padding characters
    let bytes = [];

    for (let i = 0, len = str.length; i < len; i += 4) {
        let bitString = '';
        for (let j = 0; j < 4 && i + j < len; ++j) {
            const char = str.charAt(i + j);
            const index = chars.indexOf(char);
            if (index !== -1) {
                bitString += index.toString(2).padStart(6, '0');
            }
        }

        for (let k = 0; k < bitString.length; k += 8) {
            if (k + 8 <= bitString.length) {
                const byte = bitString.substring(k, k + 8);
                bytes.push(parseInt(byte, 2));
            }
        }
    }

    return new Uint8Array(bytes);
}
export function Uint8ArrayToBase64(uint8Array: Uint8Array): string {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';
    let bitString = '';
    // Convert each byte to its 8-bit binary representation
    for (let i = 0; i < uint8Array.length; i++) {
        bitString += uint8Array[i].toString(2).padStart(8, '0');
    }

    let base64 = '';
    // Process each 6-bit segment
    for (let i = 0; i < bitString.length; i += 6) {
        const bits = bitString.substring(i, i + 6);
        // Right-pad with zeros if the last segment is less than 6 bits
        const paddedBits = bits.padEnd(6, '0');
        const index = parseInt(paddedBits, 2);
        base64 += chars[index];
    }

    // Calculate padding. Base64 output length must be divisible by 4.
    while (base64.length % 4 !== 0) {
        base64 += '=';
    }

    return base64;
}

