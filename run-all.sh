#!/bin/bash

# Configuration for Proxmox VE Connection
export PMO_PROXMOX_HOST="192.168.1.99"
export PMO_PROXMOX_PORT="8006"
export PROXMOX_USER="root"
export PROXMOX_REALM="pam"
export PROXMOX_PASSWORD="jamurkembang"

# Optional Server Config Overrides
export PMO_HOST="127.0.0.1"
export PMO_PORT="8080"
export PMO_SECRET="change-me-in-production"

# Handle graceful shutdown of backend when script exits
cleanup() {
    echo ""
    echo "Stopping background processes..."
    if [ -n "$BACKEND_PID" ]; then
        kill "$BACKEND_PID" 2>/dev/null
    fi
    exit 0
}
trap cleanup SIGINT SIGTERM EXIT

# Get project root directory
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"

# 1. Start backend in background
echo "=================================================================="
echo " Starting Proxmox Modern Orchestrator (PMO) Backend..."
echo " Target Host : https://${PMO_PROXMOX_HOST}:${PMO_PROXMOX_PORT}"
echo " User        : ${PROXMOX_USER}@${PROXMOX_REALM}"
echo "=================================================================="

cd "$PROJECT_ROOT/backend" || exit 1
go run ./cmd/server &
BACKEND_PID=$!

# Wait briefly for backend to start up
sleep 2

# 2. Start frontend in foreground
echo "=================================================================="
echo " Starting Proxmox Modern Orchestrator (PMO) Frontend..."
echo "=================================================================="

cd "$PROJECT_ROOT/frontend" || exit 1

# Install node_modules if not present
if [ ! -d "node_modules" ]; then
    echo "node_modules not found, running npm install..."
    npm install
fi

# Run vite dev server
npm run dev
