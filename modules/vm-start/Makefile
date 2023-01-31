
mkfile_path := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
FIRECRACKER ?= $(HOME)/.local/bin/firecracker-bin

start-vm:
	$(FIRECRACKER) --no-api --config-file $(mkfile_path)fc-config.json --boot-timer