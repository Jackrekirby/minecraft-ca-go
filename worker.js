// worker.js
self.importScripts('wasm_exec.js');

const go = new Go();
let wasmInstance;

self.onmessage = async (event) => {
    if (event.data.type === 'initialize') {
        const { wasmPath } = event.data;

        // Fetch and instantiate the WASM module.
        const response = await fetch(wasmPath);
        const wasmModule = await WebAssembly.instantiateStreaming(response, go.importObject);
        wasmInstance = wasmModule.instance;

        // Start the Go runtime.
        go.run(wasmInstance);

        // Notify the main thread that initialization is complete.
        self.postMessage({ type: 'initialized' });
    }

    if (event.data.type === 'runProgram') {
        if (typeof self.runProgram === 'function') {
            // Call the exported Go function `runProgram`.
            self.runProgram();
            self.postMessage({ type: 'done' });
        } else {
            console.error('runProgram is not defined in the WebAssembly module.');
            self.postMessage({ type: 'error', message: 'runProgram is not defined.' });
        }
    }

    if (event.data.type === 'onKeyDownMC') {
        if (typeof self.onKeyDownMC === 'function') {
            self.onKeyDownMC(event.data.event);
        } else {
            console.error('onKeyDownMC is not defined in the WebAssembly module.');
        }
    }
};
