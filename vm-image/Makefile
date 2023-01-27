DIRS = init kernel
ORG ?= $USER

build:
	@set -e; for d in $(DIRS); do linuxkit pkg build --org vltmn --force "$$d"; done