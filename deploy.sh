#!/bin/bash
TARGET_USER=artem
TARGET_HOST=192.168.1.118
TARGET_DIR=/home/$TARGET_USER/ein-go-client
TARGET_GOPATH=/home/$TARGET_USER/go
# That line configures the target OS as Linux, the architecture as ARM and the ARM version as 7,
# which is good for the Raspberry Pi 2 and 3 boards.
# For other versions of the Pi – A, A+, B, B+ or Zero – you’d using 6.
ARM_VERSION=6

# Executable name is assumed to be same as current directory name
EXECUTABLE=${PWD##*/}

echo "Uploading source code to Raspberry Pi..."
scp -r . $TARGET_USER@$TARGET_HOST:$TARGET_DIR

echo "Build dependencies on Raspberry Pi..."
ssh $TARGET_USER@$TARGET_HOST "cd $TARGET_DIR
go mod tidy
cd $TARGET_GOPATH/src/github.com/mcuadros/go-rpi-rgb-led-matrix/vendor/rpi-rgb-led-matrix/
git submodule update --init
make
cd $TARGET_GOPATH/src/github.com/mcuadros/go-rpi-rgb-led-matrix/
go install -v ./..."

echo "Building on Raspberry Pi..."
ssh $TARGET_USER@$TARGET_HOST "cd $TARGET_DIR && go build"
