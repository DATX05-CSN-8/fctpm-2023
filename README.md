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