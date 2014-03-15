#ifndef __INCLUDE_PROGRAMVARE_H__
#define __INCLUDE_PROGRAMVARE_H__

#define INIT    1
#define IDLE    2
#define GO_UP   3
#define GO_DOWN 4
#define E_STOP  5

// Tilstandsmaskin 
int state_machine();

// Initialiserer heisen etter oppstart. Setter heis i definert tilstand.
int initialize();

// Ser etter ordre og returnerer evt. ny tilstand.
int check_orders();

// Funksjon som kjører heisen oppover og stopper på veien hvis den skal. Returnerer tilstand.
int going_up(int *door_opened);

// Funksjon som kjører heisen nedover og stopper på veien hvis den skal. Returnerer tilstand.
int going_down(int *door_opened);

// Funksjon som sjekker om den skal stopp ift. hvilken tilstand den er i. Slukker eventuelle lys.
int stop_check(int state);

// Er i nødstopp. Sjekker om den skal fortsette der. Returnerer tilstand.
int e_stop();

#endif // #ifndef __INCLUDE_PROGRAMVARE_H__
