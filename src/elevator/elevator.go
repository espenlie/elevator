package elevator

import (
//  "fmt"
    "time"
    "drivers"
    "misc"
)
const N_FLOORS = 4
const N_BUTTONS = 3
const STOP_REVERSE_TIME = 10 // Antall Millisecond i revers ved stop

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

var Button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
    {drivers.FLOOR_UP1, drivers.FLOOR_DOWN1, drivers.FLOOR_COMMAND1},
    {drivers.FLOOR_UP2, drivers.FLOOR_DOWN2, drivers.FLOOR_COMMAND2},
    {drivers.FLOOR_UP3, drivers.FLOOR_DOWN3, drivers.FLOOR_COMMAND3},
    {drivers.FLOOR_UP4, drivers.FLOOR_DOWN4, drivers.FLOOR_COMMAND4},
}

func Elev_set_speed(speed int){
    // Hvis speed blir satt til 0, vil neste if-setning hjelpe oss til å stoppe heisen effektivt
    if (speed == 0){
        if(drivers.ReadBit(drivers.MOTORDIR)){
            drivers.ClearBit(drivers.MOTORDIR)
        }else {
            drivers.SetBit(drivers.MOTORDIR)
        }
				time.Sleep(STOP_REVERSE_TIME * time.Millisecond)
    }
    // If to start (speed > 0)
    if (speed > 0){
        drivers.ClearBit(drivers.MOTORDIR)
    }else{
        drivers.SetBit(drivers.MOTORDIR)
    }
    // Write new setting to motor.
    drivers.WriteAnalog(drivers.MOTOR, 2048 + 4*misc.Abs(speed));
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

func Elev_at_floor()bool{
    if Elev_get_floor_sensor_signal()!=-1{
        return true
    }
    return false
}

func Elev_get_floor_sensor_signal()int{
    if drivers.ReadBit(drivers.SENSOR1){
        return 1
    }else if drivers.ReadBit(drivers.SENSOR2){
        return 2
    }else if drivers.ReadBit(drivers.SENSOR3){
        return 3
    }else if drivers.ReadBit(drivers.SENSOR4){
        return 4
    }else{
        return -1
    }
}


func FloorUpdater(){
    floor := Elev_get_floor_sensor_signal()
    if (floor!=-1){
        Elev_set_floor_indicator(floor)
    }
}

func Current_floor()int{
    floor :=1
    if drivers.ReadBit(drivers.FLOOR_IND2){
        floor=floor+1}
    if drivers.ReadBit(drivers.FLOOR_IND1){
        floor=floor+2}
    return floor
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

func Elev_get_button_signal(button Elev_button, floor int)int{
    if (drivers.ReadBit(Button_channel_matrix[floor][button])){
        return 1
    }else{
        return 0
    }
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

