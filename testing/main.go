package main

import (
//  "os/exec"
    "fmt"
//  "drivers"
//  "net"
//  "misc"
//  "elevator"
    "networking"
    "misc"
//  "encoding/json"
    )

func main() {
    var message networking.Networkmessage
    message.Status = networking.Status{State: "IDLE", LastFloor:1, Source: "127.0.0.1"} 
    message.Order = networking.Order{Direction: "UP", Floor: 3, InOut: false}
    fmt.Println(message)
    networkbyte := misc.PackNetworkMessage(message)
    unpack := misc.UnpackNetworkMessage(networkbyte)
    fmt.Println(unpack.Status.State)
}

/*
func getIP() {
    cmd, err := exec.Command("ifconfig | grep 'inet addr:'").Output()
    if err != nil {
        fmt.Println(err)
    }
//  out, err := cmd.StdoutPipe()
    fmt.Printf("%s",cmd)
}


func main() {
    stat := &Status{State:      "UP",
                    LastFloor:  2,
                    Source:     "Ip"}
    mar, err := json.Marshal(stat)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(string(mar))    
    var test Order
    err = json.Unmarshal(mar, &test)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(test)
}
*/
