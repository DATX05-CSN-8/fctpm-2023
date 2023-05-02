DIRS = boottime-init shell-init kernel
ORG ?= $(USER)
LOOP ?= /dev/loop0
IMG_NAME ?= fc-image
OUTP_IMAGE ?= fc-image
INIT_NAME ?= boottime-init

clean: 
	rm -rf out \
	rm -rf */PKG_TAG
	rm -rf img/*.yml

build: out/fc-disk.img

out: 
	mkdir -p out

out/build-pkg: out
	@set -e; for d in $(DIRS); do \
		linuxkit pkg build --org $(ORG) --force "$$d" && \
		linuxkit pkg show-tag --org $(ORG) "$$d" > $$d/PKG_TAG; \
	done
	touch out/build-pkg

img/$(IMG_NAME).yml: out/build-pkg
	KERNEL_IMAGE=$$(cat kernel/PKG_TAG) INIT_IMAGE=$$(cat $(INIT_NAME)/PKG_TAG) envsubst < img/$(IMG_NAME).yml.tmpl > img/$(IMG_NAME).yml

out/$(OUTP_IMAGE)-kernel: img/$(IMG_NAME).yml
	linuxkit build --format kernel+initrd --decompress-kernel -dir out -name $(OUTP_IMAGE) img/$(IMG_NAME).yml

out/fc-disk.img: out/fc-image-kernel
	linuxkit build --format tar img/$(IMG_NAME).yml && \
	sudo losetup -d $(LOOP) || true
	dd if=/dev/zero of=$(@F) bs=100M count=10
	sudo mkfs.ext4 $(@F)
	sudo losetup $(LOOP) $(@F)
	mkdir fc.img
	sudo mount $(LOOP) fc.img
	sudo tar xf $(IMG_NAME).tar -C fc.img
	sudo umount fc.img
	sudo losetup -d $(LOOP)
	rmdir fc.img
	mv $(@F) $@
	rm $(IMG_NAME).tar
