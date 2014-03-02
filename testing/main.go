package main

import (
    "os/exec"
    "fmt"
//  "fmt"
//  "drivers"
//  "net"
//  "misc"
//  "elevator"
    "encoding/json"
    )

type Order struct {
    Direction   string
    Floor       int
}

type Status struct {
    State       string
    LastFloor   int
    Source      string
}

func getIP() {
    cmd, err := exec.Command("ifconfig | grep 'inet addr:'").Output()
    if err != nil {
        fmt.Println(err)
    }
//  out, err := cmd.StdoutPipe()
    fmt.Printf("%s",cmd)
}
func main() {
//  getIP()
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
