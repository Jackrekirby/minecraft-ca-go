# Set environment variables for WebAssembly
$env:GOOS = "js"
$env:GOARCH = "wasm"

# Build the Go project to WebAssembly
go build -o main.wasm
# go build -o main.wasm -ldflags="-s -w" -gcflags=all=-l


# Check if go build was successful
if ($?) {
    Write-Output "built main.wasm"
}
else {
    Write-Output "Error: go build failed. Stopping the script."
    exit 1  # Exit the script with a non-zero status
}

# Get GOROOT to find wasm_exec.js
$goroot = & go env GOROOT

# Path to wasm_exec.js
$wasmExecPath = Join-Path $goroot "misc\wasm\wasm_exec.js"

# Check if wasm_exec.js exists and copy it to the current directory
if (Test-Path $wasmExecPath) {
    Copy-Item -Path $wasmExecPath -Destination . -Force
    Write-Output "wasm_exec.js has been copied to the current directory. `"$wasmExecPath`""
}
else {
    Write-Output "Error: wasm_exec.js was not found in GOROOT. Please check your Go installation."
}
