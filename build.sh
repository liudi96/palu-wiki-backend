#!/bin/bash

echo "Building Go backend application..."
go build -o server_app cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "Build successful: server_app"
else
    echo "Build failed!"
    exit 1
fi
