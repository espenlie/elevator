#include <stdio.h>
#include <time.h>
#include "channels.h"
#include "io.h"
#include "heis.h"
#include "elev.h"


void refresh_floor_indicator(){
    int floor = elev_get_floor_sensor_signal();
    if(floor != -1)
       elev_set_floor_indicator(floor);
    return;
}


int current_floor(){
    int floor=0;
    if (io_read_bit(FLOOR_IND2))
        floor++;
    if (io_read_bit(FLOOR_IND1))
        floor=floor+2;
    return floor;
}


int door_open(){
    return io_read_bit(DOOR_OPEN);
}


void close_door(int *door_opened){
    int end_t = time(0);
    int difftime = end_t-(*door_opened);
    if (difftime >= OBSTRUCTION_TIME)
        elev_set_door_open_lamp(0);
}

