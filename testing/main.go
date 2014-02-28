package main

import (
    "fmt"
    "misc"
//  "elevator"
)

func main() {
    config := misc.LoadConfig("/home/student/LL/elevator/config/conf.json")
    fmt.Println(config.Default_Listen_Port)
    fmt.Println(config.Default_Dial_Port)
//  for _,key := range config.Elevators {
//      fmt.Println(key.Address)
//  }
    fmt.Println(config.Elevators)
}
