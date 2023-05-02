#include <unistd.h>
#include <fcntl.h>
#include <sys/mman.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <sys/reboot.h> 
#include <stdio.h>

/*
 * This little program is a mock 'init' process. It writes the firecracker debug port.
 * It will write a timestamp to the log.
 */

#define MAGIC_MMIO_SIGNAL_GUEST_BOOT_COMPLETE 0xd0000000
#define MAGIC_VALUE_SIGNAL_GUEST_BOOT_COMPLETE 123

int main(void)
{
    // set up boot timer device mmio
    int fd = open("/dev/mem", (O_RDWR | O_SYNC | O_CLOEXEC));
    if (fd == -1) {
        perror("Error opening file: ");
        return -1;
    }
    
    int mapped_size = getpagesize();
    char *map_base = mmap(NULL,
            mapped_size,
            PROT_WRITE,
            MAP_SHARED,
            fd,
            MAGIC_MMIO_SIGNAL_GUEST_BOOT_COMPLETE);

    // write guest boot complete command
    *map_base = MAGIC_VALUE_SIGNAL_GUEST_BOOT_COMPLETE;
    msync(map_base, mapped_size, MS_ASYNC);
    return 0;
}
