#include <sys/io.h>
#include <stdio.h>

/*
 * This little program is a mock 'init' process. It writes the firecracker debug port.
 * It will write a timestamp to the log.
 */

#define PORT_FC     0x03f0
#define PORT_FC_VAL 123

int main(void)
{
    int r;

    r = ioperm(PORT_FC, 1, 1);
    if (r) {
        fprintf(stderr, "Error setting up port access to 0x%x, quitting\n", PORT_FC);
        return -1;
    }

    outb(PORT_FC_VAL, PORT_FC);
    return 0;
}
