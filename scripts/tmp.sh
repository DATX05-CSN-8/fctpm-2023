qemu-system-x86_64 -L ../bin/bios -smp 1 -m 256 -accel kvm -cpu host -machine pc -display none -nographic -vga none \
    -initrd ./out/with-init-initrd.img \
    -no-reboot --no-acpi \
    -chardev socket,id=chrtpm,path=/tmp/swtpm-go/24a780a3-2035-4d40-b246-75749c5e85b1/socket \
    -tpmdev emulator,id=tpm0,chardev=chrtpm \
    -device tpm-tis,tpmdev=tpm0 \
    -nic none \
    -kernel ./out/with-init-kernel -append "console=ttyS0 panic=-1 tpm_tis.force=1 ima=on ima_policy=tcb lsm=integrity"


qemu-system-x86_64 -L ../bin/bios -smp 1 -m 256 -accel kvm -cpu host -machine pc -display none -nographic -vga none \
    -initrd ./out/with-init-initrd.img \
    -no-reboot --no-acpi \
    -nic none \
    -kernel ./out/with-init-kernel -append "console=ttyS0 panic=-1"

/home/melker/firecracker/build/cargo_target/x86_64-unknown-linux-musl/debug/firecracker \
    --no-api --config-file /tmp/firecracker-config/simple.json

/home/melker/fctpm-2023/modules/firecracker/bin/firecracker \
    --no-api --config-file /tmp/firecracker-config/simple.json

/home/melker/cloud-hypervisor/target/release/cloud-hypervisor \
    --kernel ./vm-image/out/with-init-kernel \
    --initramfs ./vm-image/out/with-init-initrd.img \
    --cmdline "console=ttyS0 panic=-1 tpm_tis.force=1 ima=on ima_policy=tcb lsm=integrity" \
    --cpus boot=2 --memory size=256 \
    --tpm socket=/tmp/swtpm-go/a02c4a6b-59be-4db5-af08-a47b0351a9ec/socket