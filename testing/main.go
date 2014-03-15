package main

import (
    "fmt"
    "misc"
    )
/*
func main() {
    var message networking.Networkmessage
    message.Status = networking.Status{State: "IDLE", LastFloor:1, Source: "127.0.0.1"} 
    message.Order = networking.Order{Direction: "UP", Floor: 3, InOut: false}
    fmt.Println(message)
    networkbyte := misc.PackNetworkMessage(message)
    unpack := misc.UnpackNetworkMessage(networkbyte)
    fmt.Println(unpack.Status.State)
}
*/


func main() {
    fmt.Println(misc.GetLocalIP())
}
