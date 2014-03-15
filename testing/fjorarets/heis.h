#ifndef __INCLUDE_HEIS_H__ 
#define __INCLUDE_HEIS_H__

#define SPEED 500   // Definerer motorfart    
#define OBSTRUCTION_TIME 3   // Definerer obstrusjontiden
#define STOP_REVERSE_TIME 4230 // Definerer motorens reverstid ved stop (µs)
        
// Sjekker etasjesensorer og setter etasjeindikatorene etter hvor heisen er.
void refresh_floor_indicator();

// Returnerer forrige passerte etasje.
int current_floor();

// Returnerer om døren er open eller ikke.
int door_open();

// Lukker døren når/hvis OBSTRUCTION_TIME er nådd.
void close_door(int *door_opened);


#endif // #ifndef __INCLUDE_HEIS_H__

