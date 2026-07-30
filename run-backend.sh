#!/bin/bash

# Configuration for Proxmox VE Connection
export PMO_PROXMOX_HOST="192.168.1.99"
export PMO_PROXMOX_PORT="8006"
export PROXMOX_USER="root"
export PROXMOX_REALM="pam"
export PROXMOX_PASSWORD="jamurkembang"

# Optional Server Config Overrides
export PMO_HOST="0.0.0.0"
export PMO_PORT="8080"
export PMO_SECRET="change-me-in-production"

echo "=================================================================="
echo " Starting Proxmox Modern Orchestrator (PMO) Backend"
echo " Target Host : https://${PMO_PROXMOX_HOST}:${PMO_PROXMOX_PORT}"
echo " User        : ${PROXMOX_USER}@${PROXMOX_REALM}"
echo " Listen Addr : http://${PMO_HOST}:${PMO_PORT}"
echo "=================================================================="

# Navigate to the backend folder and run the Go server
cd "$(dirname "$0")/backend" || exit 1
go run ./cmd/server
