#include <stdio.h>
#include "elev.h"
#include "io.h"
#include "heis.h"
#include "brukergrensesnitt.h"


void set_button_lights(int state){
    int floor, button_type;
    for (floor = 0 ; floor < 4 ; floor++){
        for (button_type = 0 ; button_type < 3 ; button_type++){
            if((floor == 0 && button_type == 1) || (floor == 3 && button_type == 0)){
            }
            else if (elev_get_button_signal(button_type, floor)){
                elev_set_button_lamp(button_type, floor, 1);
            }
        }
    }
}


int check_estop(){
    if(elev_get_stop_signal()){
        elev_set_stop_lamp(1);
        reset_all_orders();
        elev_set_speed(0);
        return 1;
    }
    return 0;
}


int any_order_above(){
    int floor, button_type;
    for (floor = current_floor() +1 ; floor < 4 ; floor++){
        for (button_type = 0 ; button_type < 3 ; button_type++ ){
            if(!((floor == 0 && button_type == 1) || (floor == 3 && button_type == 0))){
                if (elev_get_button_lamp_signal(button_type, floor))
                    return 1; 
            }
        }
    }
    return 0;
}


int any_order_below(){
    int floor, button_type;
    for (floor = 0 ; floor < current_floor() ; floor++){
        for (button_type = 0 ; button_type < 3 ; button_type++ ){
            if(!((floor == 0 && button_type == 1) || (floor == 3 && button_type == 0))){
                if (elev_get_button_lamp_signal(button_type, floor))
                    return 1;
            } 
        }
    }
    return 0;
}


void reset_all_orders(){
    int floor, button_type;
    for (floor = 0 ; floor < N_FLOORS ; floor++){
        for (button_type = 0 ; button_type < 3 ; button_type++ ){
            if(!((floor == 0 && button_type == 1) || (floor == 3 && button_type == 0))){
                elev_set_button_lamp(button_type, floor, 0);
            }
        }
    }
}
