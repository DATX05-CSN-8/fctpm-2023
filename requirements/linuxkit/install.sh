#!/bin/bash

curl -sSL -O https://github.com/linuxkit/linuxkit/releases/download/v1.0.1/linuxkit-linux-amd64
sudo cp linuxkit-linux-amd64 /usr/local/bin/linuxkit
rm -rf linuxkit-linux-amd64
sudo chmod +x /usr/local/bin/linuxkit