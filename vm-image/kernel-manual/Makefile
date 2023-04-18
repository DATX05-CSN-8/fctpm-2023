CONFIG_NAME = virtio-tpm-enabled
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
KERNEL_DIR ?= $(HOME)/linux

clean_config:
	rm -rf $(KERNEL_DIR)/.config

$(KERNEL_DIR)/.config: $(CONFIG_NAME).config
	cp $(CONFIG_NAME).config $(KERNEL_DIR)/.config

config: $(KERNEL_DIR)/.config
	make -C $(KERNEL_DIR) oldconfig

out:
	mkdir -p out

clean:
	rm -rf out

deep_clean: clean clean_config
	make -C $(KERNEL_DIR) clean

build $(KERNEL_DIR)/arch/x86_64/boot/bzImage out/modules: out config
	make -C $(KERNEL_DIR) -j "$$(getconf _NPROCESSORS_ONLN)" KCFLAGS="-fno-pie"
	make -C $(KERNEL_DIR) -j "$$(getconf _NPROCESSORS_ONLN)" KCFLAGS="-fno-pie" \
		INSTALL_MOD_PATH=$(ROOT_DIR)/out/modules modules_install

package out/kernel.tar out/kernel: build
	cp $(KERNEL_DIR)/arch/x86_64/boot/bzImage out/kernel
	bsdtar cf out/kernel.tar out/modules
	rm -rf out/modules