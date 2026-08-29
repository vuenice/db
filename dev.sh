#!/bin/bash

echo "Starting backend..."
(cd backend && go run ./cmd/chatdb) &
BACKEND_PID=$!
echo "Waiting for backend to start..."
sleep 5

echo "Starting frontend..."
(cd frontend && npm run dev) &
FRONTEND_PID=$!

echo "Development servers started! Press Ctrl+C to stop both."

# Trap Ctrl+C (SIGINT) to kill both background processes
trap "kill $BACKEND_PID $FRONTEND_PID; exit" SIGINT

# Wait for background processes to keep script running
wait
