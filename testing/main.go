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

func getIP() {
    oneliner := "ifconfig | grep 129.241.187 | cut -d' ' -f2 | cut -d':' -f1"
    cmd, err := exec.Command("bash",
    if err != nil {
        fmt.Println(err)
    }
    fmt.Printf("%s",cmd)
}

/*
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
