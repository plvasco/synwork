#!/bin/bash

# Set variables
MOUNT_POINT="/home/Docuements/share"
SHARE="//192.168.1.100/share"
USERNAME="your_user"
PASSWORD="your_pass"

# 1. Create the mount directory if it doesn't exist
if [ ! -d "$MOUNT_POINT" ]; then
    echo "Creating mount directory at $MOUNT_POINT..."
    sudo mkdir -p "$MOUNT_POINT"
fi

# 2. Install cifs-utils if not already installed
if ! dpkg -s cifs-utils &> /dev/null; then
    echo "Installing cifs-utils..."
    sudo apt update && sudo apt install -y cifs-utils
else
    echo "cifs-utils already installed."
fi

# 3. Mount the CIFS share
echo "Mounting CIFS share..."
sudo mount -t cifs "$SHARE" "$MOUNT_POINT" -o username="$USERNAME",password="$PASSWORD"

# Check if mount was successful
if mountpoint -q "$MOUNT_POINT"; then
    echo "Mount successful at $MOUNT_POINT."
else
    echo "Failed to mount CIFS share."
    exit 1
fi
