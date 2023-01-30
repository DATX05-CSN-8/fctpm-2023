#include <sys/io.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>

/*
 * This little program is a mock 'init' process. It writes the firecracker debug port.
 * It will write a timestamp to the log.
 */

#define PORT_FC     0x03f0
#define PORT_FC_VAL 123
#define INP_LEN 16

int main(void)
{
    int r;
    char inp[INP_LEN];

    while(1) {
        if (fgets(inp, INP_LEN, stdin) == NULL) {
            fprintf(stderr, "Error getting input data\n.");
            exit(0);
        }
        inp[strcspn(inp, "\n")] = 0;
        if (!strcmp(inp, "exit"))
            break;
    }
    r = ioperm(PORT_FC, 1, 1);
    if (r) {
        fprintf(stderr, "Error setting up port access to 0x%x, quitting\n", PORT_FC);
        return -1;
    }

    outb(PORT_FC_VAL, PORT_FC);
    return 0;
}
