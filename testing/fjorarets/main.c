#include <stdio.h>
#include "tilstandsmaskin.h"
#include "elev.h"


int main()
    {
    // Initialize hardware
    if (!elev_init()) {
        printf(__FILE__ ": Unable to initialize elevator hardware\n");
        return 1;
    }

    //Starter heis
    state_machine();
    return 0;
}

