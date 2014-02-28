package main

import (
    ."fmt"
    ."net"
    "time"
    "strconv"
)

const (
    Single = iota
    Multi
)

type sendMessage_s struct {
    sendType    int
    msg         string
}


func main(){
    initConn, err := Dial("tcp", "localhost:6969")
    if err != nil {
        Println(err)
    }
    var counter int
    var sendMessage string
    for {
        sendMessage=strconv.Itoa(counter)+"EOL"
        initConn.Write(append([]byte(sendMessage), []byte{0}...))
        counter=counter+1
        time.Sleep(1000 * time.Millisecond)
    }
}
