package main

import (
    "os/exec"
    "fmt"
    )

func getIP() {
    cmd, err := exec.Command("ifconfig | grep 'inet addr:'").Output()
    if err != nil {
        fmt.Println(err)
    }
//  out, err := cmd.StdoutPipe()
    fmt.Printf("%s",cmd)
}

//import (
//  "fmt"
//  "drivers"
//  "net"
//  "misc"
//  "elevator"
//)

//const SPEED1 = 4024
//func main() {
    
//  drivers.IoInit()
//  for{
//      if drivers.ReadBit(drivers.SENSOR1){
//          drivers.ClearBit(drivers.MOTORDIR)
//          drivers.WriteAnalog(drivers.MOTOR,SPEED1)
//      }
//      if drivers.ReadBit(drivers.SENSOR4){
//          drivers.SetBit(drivers.MOTORDIR)
//          drivers.WriteAnalog(drivers.MOTOR,SPEED1)
//      }
//      if drivers.ReadBit(drivers.STOP){
//          drivers.SetBit(drivers.MOTORDIR)
//          drivers.WriteAnalog(drivers.MOTOR,0)
//      }
//  }
//}
func main() {
    getIP()
}
