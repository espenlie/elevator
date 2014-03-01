package elevator

import (
//  ."fmt"
//  "time"
    "drivers"
)

func FloorUpdater(){
    for{
        if drivers.ReadBit(drivers.SENSOR1){
            drivers.ClearBit(drivers.FLOOR_IND1)
            drivers.ClearBit(drivers.FLOOR_IND2)
        }
        if drivers.ReadBit(drivers.SENSOR2){
            drivers.ClearBit(drivers.FLOOR_IND1)
            drivers.SetBit(drivers.FLOOR_IND2)
        }
        if drivers.ReadBit(drivers.SENSOR3){
            drivers.SetBit(drivers.FLOOR_IND1)
            drivers.ClearBit(drivers.FLOOR_IND2)
        }
        if drivers.ReadBit(drivers.SENSOR4){
            drivers.SetBit(drivers.FLOOR_IND1)
            drivers.SetBit(drivers.FLOOR_IND2)
        }
    }
}
