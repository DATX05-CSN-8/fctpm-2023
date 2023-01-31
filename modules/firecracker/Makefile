
REL_DIR=firecracker/build/cargo_target/x86_64-unknown-linux-musl/release

build: bin/firecracker

bin/firecracker: $(REL_DIR)/firecracker
	mkdir -p $(@D)
	cp $< $(@D)

$(REL_DIR)/firecracker: firecracker
	cd $< && git am ../patches/*
	cd $< && tools/devtool -y build --release

firecracker:
	git clone https://github.com/DATX05-CSN-8/firecracker.git

reset: firecracker
	cd $< && git reset --hard origin/main && git am --abort


clean:
	rm -rf bin $(REL_DIR)/firecracker

clean-git:
	rm -rf firecracker