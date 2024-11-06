# Set environment variables for Windows (adjust GOARCH as needed)
$env:GOOS = "windows"
$env:GOARCH = "amd64"  # Or use "386" for 32-bit

# Build the Go project for Windows
go build -o main.exe
