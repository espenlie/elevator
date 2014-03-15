#include <stdio.h>
#include <time.h>
#include "tilstandsmaskin.h"
#include "brukergrensesnitt.h"
#include "heis.h"
#include "elev.h"


int state_machine(){
    int state = INIT;
    int door_opened;
    while(1){
       switch (state) {
            case INIT:
                state = initialize();
                break;
            case IDLE:
                close_door(&door_opened);
                state = check_orders();
                break;
            case GO_UP:
                state = going_up(&door_opened);
                break;
            case GO_DOWN:
                state = going_down(&door_opened);
                break;
            case E_STOP:
                state = e_stop(); 
                break;
        }
        refresh_floor_indicator();
        if (state != INIT){
            set_button_lights();
            if (check_estop())
                state = E_STOP;
        }
    }
    return 1;
}


int initialize(){
    elev_set_speed(-SPEED);
    if(elev_get_floor_sensor_signal() != -1){
        elev_set_speed(0);
        return IDLE;
    }
    return INIT;
}


int check_orders(){
    if (any_order_below())
        return GO_DOWN;
    if (current_floor() != 0){
        if (elev_get_button_lamp_signal(BUTTON_CALL_DOWN, current_floor()))
            return GO_DOWN;
    }
    if (any_order_above())
        return GO_UP;
    if (current_floor() != 3){
        if (elev_get_button_lamp_signal(BUTTON_CALL_UP, current_floor()))
            return GO_UP;
    }
    if (elev_get_button_lamp_signal(BUTTON_COMMAND, current_floor()))
        return GO_UP;
    return IDLE;
}


int going_up (int *door_opened){
    if(stop_check(GO_UP)){
        *door_opened = time(0); 
        if (!door_open())
            elev_set_speed(0);
        elev_set_door_open_lamp(1);
    }
    else if (!any_order_above()){
        return IDLE;
    }
    int end_t = time(0);
    int difftime = end_t - (*door_opened);
    if (difftime >= OBSTRUCTION_TIME){
        elev_set_speed(SPEED);
        elev_set_door_open_lamp(0);
    }
    return GO_UP;
}


int going_down (int *door_opened){
    if (stop_check(GO_DOWN)){
        *door_opened=time(0); 
        if (!door_open())
            elev_set_speed(0);
        elev_set_door_open_lamp(1);
    } 
    else if (!any_order_below()){
        return IDLE;
    }
    int end_t = time(0);
    int difftime = end_t - (*door_opened);
    if (difftime >= OBSTRUCTION_TIME){
        elev_set_speed(-SPEED);
        elev_set_door_open_lamp(0);
    }
    return GO_DOWN;
}


int e_stop(){
    int floor;
    for (floor = 0 ; floor < N_FLOORS ; floor++){
        if (elev_get_button_lamp_signal(BUTTON_COMMAND, floor)){
            elev_set_stop_lamp(0);
            return IDLE;
        }
    }
    return E_STOP;
}


int stop_check (int state){
    int floor = elev_get_floor_sensor_signal();
    if (elev_get_obstruction_signal())
        return 1;
    if (floor == -1)
        return 0;
    if (state == GO_UP){
        if (floor == 3){
            if (elev_get_button_lamp_signal(BUTTON_COMMAND, floor)){
                elev_set_button_lamp(BUTTON_COMMAND, floor, 0);
                return 1;
            }
        }
        else if (elev_get_button_lamp_signal(BUTTON_CALL_UP, floor) || elev_get_button_lamp_signal(BUTTON_COMMAND, floor)){
            elev_set_button_lamp(BUTTON_COMMAND, floor, 0);
            elev_set_button_lamp(BUTTON_CALL_UP, floor, 0);
            return 1;
        }
    }
    else if (state == GO_DOWN){
        if (floor == 0){
            if (elev_get_button_lamp_signal(BUTTON_COMMAND, floor)){
                elev_set_button_lamp(BUTTON_COMMAND, floor, 0);  
                return 1;
            }          
        }
        else if (elev_get_button_lamp_signal(BUTTON_CALL_DOWN, floor) || elev_get_button_lamp_signal(BUTTON_COMMAND, floor)){
            elev_set_button_lamp(BUTTON_COMMAND, floor, 0);
            elev_set_button_lamp(BUTTON_CALL_DOWN, floor, 0);
            return 1;
        }
    }
    return 0;
}
