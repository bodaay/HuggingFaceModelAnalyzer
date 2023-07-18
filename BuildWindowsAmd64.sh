#!/bin/bash

output_dir="output"
binaryName="hfanalyzer"
read -r version < VERSION

echo "Building for Windows/amd64..."
GOOS=windows GOARCH=amd64 go build -o "$output_dir/${binaryName}_windows_amd64_$version.exe" main.go
echo "Build completed. Binary output: $output_dir/${binaryName}_windows_amd64_$version.exe"
