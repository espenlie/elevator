package elevator

import (
//  "fmt"
    "time"
    "drivers"
)
const N_FLOORS = 4
const N_BUTTONS = 3

type Elev_button int

const(
    BUTTON_CALL_UP Elev_button = iota
    BUTTON_CALL_DOWN
    BUTTON_COMMAND
)

var Lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{
    {drivers.LIGHT_UP1, drivers.LIGHT_DOWN1, drivers.LIGHT_COMMAND1},
    {drivers.LIGHT_UP2, drivers.LIGHT_DOWN2, drivers.LIGHT_COMMAND2},
    {drivers.LIGHT_UP3, drivers.LIGHT_DOWN3, drivers.LIGHT_COMMAND3},
    {drivers.LIGHT_UP4, drivers.LIGHT_DOWN4, drivers.LIGHT_COMMAND4},
}

var button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
    {drivers.FLOOR_UP1, drivers.FLOOR_DOWN1, drivers.FLOOR_COMMAND1},
    {drivers.FLOOR_UP2, drivers.FLOOR_DOWN2, drivers.FLOOR_COMMAND2},
    {drivers.FLOOR_UP3, drivers.FLOOR_DOWN3, drivers.FLOOR_COMMAND3},
    {drivers.FLOOR_UP4, drivers.FLOOR_DOWN4, drivers.FLOOR_COMMAND4},
}

func Elev_set_floor_indicator(floor int){
    if (floor == 3 || floor == 4){
        drivers.SetBit(drivers.FLOOR_IND1)
    }else{
        drivers.ClearBit(drivers.FLOOR_IND1)}
    if (floor == 2 || floor ==4){
        drivers.SetBit(drivers.FLOOR_IND2)
    }else{
        drivers.ClearBit(drivers.FLOOR_IND2)}
}

func FloorUpdater(){
    for{
        if drivers.ReadBit(drivers.SENSOR1){
            Elev_set_floor_indicator(1)
        }
        if drivers.ReadBit(drivers.SENSOR2){
            Elev_set_floor_indicator(2)
        }
        if drivers.ReadBit(drivers.SENSOR3){
            Elev_set_floor_indicator(3)
        }
        if drivers.ReadBit(drivers.SENSOR4){
            Elev_set_floor_indicator(4)
        }
        time.Sleep(100 * time.Millisecond)
    }
}

func Elev_set_button_lamp(button Elev_button, floor int, value int){
    if (value == 1){
        drivers.SetBit(Lamp_channel_matrix[floor-1][button])
    }else{
        drivers.ClearBit(Lamp_channel_matrix[floor-1][button])
        }
}

func Elev_set_door_open_lamp(value int){
    if (value == 1){
        drivers.SetBit(drivers.DOOR_OPEN)
        }else{
        drivers.ClearBit(drivers.DOOR_OPEN)}
}

func Elev_set_stop_lamp(value int){
    if (value==1){
        drivers.SetBit(drivers.LIGHT_STOP)
    }else{
        drivers.ClearBit(drivers.LIGHT_STOP)}
}


func Elev_init(){
    // Zero all floor button lamps
    for i := 1; i < N_FLOORS+1; i++ {
        if (i != 1){
            Elev_set_button_lamp(BUTTON_CALL_DOWN, i, 0)
        }
        if (i != N_FLOORS){
            Elev_set_button_lamp(BUTTON_CALL_UP, i, 0)
        }
        Elev_set_button_lamp(BUTTON_COMMAND, i, 0)
    }

    // Clear stop lamp, door open lamp, and set floor indicator to ground floor.
    Elev_set_stop_lamp(0);
    Elev_set_door_open_lamp(0);
    Elev_set_floor_indicator(0);
}

