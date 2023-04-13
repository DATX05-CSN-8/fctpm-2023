DIRS = init simple_init kernel-manual
ORG ?= $(USER)
LOOP ?= /dev/loop0
IMG_NAME = fc-image

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
	KERNEL_IMAGE=$$(cat kernel-manual/PKG_TAG) INIT_IMAGE=$$(cat simple_init/PKG_TAG) envsubst < img/$(IMG_NAME).yml.tmpl > img/$(IMG_NAME).yml

out/fc-image-kernel: img/$(IMG_NAME).yml
	linuxkit build --format kernel+initrd --decompress-kernel -dir out img/$(IMG_NAME).yml

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
