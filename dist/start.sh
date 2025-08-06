#!/bin/bash
# Mosaic App Startup Script

# Load environment variables
if [ -f ".env" ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Set default port if not set
export SERVER_PORT=${SERVER_PORT:-8080}

echo "Starting Mosaic App on port $SERVER_PORT..."
./mosaic
