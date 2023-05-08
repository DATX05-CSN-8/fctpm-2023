
REL_DIR=firecracker/build/cargo_target/x86_64-unknown-linux-musl/release

build: bin/firecracker

bin/firecracker: $(REL_DIR)/firecracker
	mkdir -p $(@D)
	cp $< $(@D)

$(REL_DIR)/firecracker: firecracker
	for f in ./patches/*; do patch -d $< -p1 -s --forward < $$f; done
	cd $< && tools/devtool -y build --release

firecracker:
	git clone https://github.com/firecracker-microvm/firecracker.git
	cd firecracker && git reset --hard c2d11ff24c397b1c516b5b10b998bfd4860a7ed5

clean:
	rm -rf bin $(REL_DIR)/firecracker firecracker