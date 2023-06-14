# Trust in Lightweight Virtual Machines: Integrating TPMs into Firecracker
This repository contains the source code to reproduce the experiments made in a Master's thesis project at Chalmers University of Technology (CTH), Sweden.
The thesis project was done by Alexandra Parkegren and Melker Veltman and supervised by Victor Morel.

## Abstract
When software services use cloud providers to run their workloads, they place implicit trust in the cloud provider, without an explicit trust relationship.
One way to achieve such explicit trust in a computer system is to use a hardware Trusted Platform Module (TPM), a coprocessor for trusted computing. 
However, in the case of managed platform-as-a-service (PaaS) offerings, there is currently no cloud provider that exposes TPM capabilities. 
In this paper, we improve trust by integrating a virtual TPM device into the Firecracker hypervisor, originally developed by Amazon Web Services. 
In addition to this, multiple performance tests along with an attack surface analysis are performed to evaluate the impact of the changes introduced. 

We discuss the results and conclude that the slight performance decrease and attack surface increase are acceptable trade-offs in order to enable trusted computing in PaaS offerings. 
## Getting Started
### Prerequisites
The instructions are made for Ubuntu 22.04 running on an x86-64 processor capable of virtualisation with at least 8GB memory. 
It is assumed that the working directory is the root of the project repository. Furthermore, some packages need to be installed and the current user need to be added to the kvm user group to be able to interact with KVM. 
To accomplish these installations, the following commands should be run:
```bash
$ sudo apt-get update
$ sudo apt-get install docker.io make patch
$ sudo usermod -a -G kvm $USER
$ sudo usermod -a -G docker $USER
```
Also, Linuxkit and swtpm needs to be installed by running the following command:
```bash
$ make -C requirements
$ swtpm --version # Prints the swtpm version if it is available
$ linuxkit version # Prints the linuxkit version if it is available
```
Since the group assignment of the current user has been changed, the terminal session needs to be restarted before the next step.

### Compiling Firecracker and Building a Linux Kernel
To be able to run Firecracker with the added TPM functionality, it needs to be compiled from source with the added patches. 
For convenience, the changes needed are available as patch files in the repository of the project, and the binary can be compiled using the following command:
```bash
$ make -C modules/firecracker build
$ # Creates a firecracker binary in modules/firecracker/bin/firecracker
```
In order to start a Firecracker VM, a Linux kernel and accompanying init program is also needed. 
These can be built using the following command:
```bash
$ INIT_NAME=shell-init make -C vm-image out/fc-image-kernel $ # Creates kernel and init program at vm-image/out
```
### Starting a vTPM and Running a VM
As the Firecracker VM needs the TPM to be available once it starts, the swtpm process needs to be started before the VM. 
The commands below performs the key generation process for the TPM and then starts the swtpm process with a specified UNIX socket:
```bash
$ mkdir -p .tmp/tpm
$ swtpm_setup --tpm-state .tmp/tpm \
    --createek --tpm2 --create-ek-cert \
    --create-platform-cert --lock-nvram
$ swtpm socket --tpmstate dir=.tmp/tpm --tpm2 \
    --ctrl type=unixio,path=.tmp/tpm/socket --flags startup-clear
```
As this process needs to be kept running when the VM is started, another terminal
is needed to start the Firecracker process. 
The command to run is the following:
```bash
$ modules/firecracker/bin/firecracker --no-api \
    --config-file modules/vm-start/shell-config.json
```
If successful, it prints the Linux kernel output and returns the user to a shell inside the VM. 
Running the following command displays that the TPM is available from within the VM.
```bash
$ ls -l /dev/tpm0
```
The VM can be stopped by running exit from within the VM shell. 
After that, the swtpm process can be stopped by pressing ctrl+C.

