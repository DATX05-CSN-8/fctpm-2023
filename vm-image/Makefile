DIRS = init kernel
ORG ?= $(USER)

build-pkg:
	@set -e; for d in $(DIRS); do \
		linuxkit pkg build --org $(ORG) --force "$$d" && \
		linuxkit pkg show-tag --org $(ORG) "$$d" > $$d/PKG_TAG; \
		done

build-img:
	KERNEL_IMAGE=$$(cat kernel/PKG_TAG) INIT_IMAGE=$$(cat init/PKG_TAG) envsubst < img/fc-image.yml.tmpl > img/fc-image.yml && \
	mkdir -p out && \
	linuxkit build --format kernel+initrd --decompress-kernel -dir out img/fc-image.yml