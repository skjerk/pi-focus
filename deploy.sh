#!/bin/bash
TARGET_USER=pi
TARGET_HOST=octopi.local
TARGET_DIR=
ARM_VERSION=7
 
# Executable name is assumed to be same as current directory name
# Other files or direcories are specified in OTHER
EXECUTABLE=${PWD##*/} 
OTHER="templates public"
 
echo "Building for Raspberry Pi..."
env GOOS=linux GOARCH=arm GOARM=$ARM_VERSION go build
 
echo "Uploading to Raspberry Pi..."
scp -r $EXECUTABLE $OTHER $TARGET_USER@$TARGET_HOST:$TARGET_DIR