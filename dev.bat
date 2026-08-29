@echo off
echo Starting backend...
start cmd /k "cd backend && go run ./cmd/chatdb"

echo Starting frontend...
start cmd /k "cd frontend && npm run dev"

echo Development servers started in separate windows!
