package main

import (
    ."fmt"
    ."net"
    ."strings"
//  "os/exec"
//  "time"
//  "strconv"
)

func receiver(connection Conn, received chan string){
    var buf [1024]byte
    for {
        _, err := connection.Read(buf[0:])
        if err != nil {
//          status <- "YO"
            Println(err)
            return
        }
        received <- Split(string(buf[0:]), "EOL")[0]
    }
}


func listener(conn *TCPListener, newConn_c chan Conn){
    for {
        newConn, err := conn.Accept()
        if err != nil {
            Println(err)
        }
        newConn_c <- newConn
    }
}

func networking(newConn_c chan Conn, generatedMsgs_c chan string, receivedMsgs_c chan string, dialConn_c chan Conn) {
    var newConn Conn
    for{
        select {
        case newConn = <- dialConn_c:
		case sendMsg := <- generatedMsgs_c:
			newConn.Write(append([]byte(sendMsg), []byte{0}...))
        case msg := <-receivedMsgs_c:
            Println(msg)
        case newConn := <- newConn_c:
//			connMap[newConn.LocalAddr().String()] = newConn
			go receiver(newConn, receivedMsgs_c)
	
        }
    }
}

func dialer(elevators string, port string, dialconn_c chan Conn){
	for {
			dialConn, err := Dial("tcp", elevator+port)
			if err != nil {
				Println(err)
			}else{
                dialconn_c <-dialConn
                break
            }
		}
	}

func main(){

    listenAddr, _ := ResolveTCPAddr("tcp", ":6969")
    listenConn, _ := ListenTCP("tcp", listenAddr)
    receivedMsgs_c  := make(chan string)
    generatedMsgs_c  := make(chan string)
    newConn_c       := make(chan Conn, 10)
    dialconn_c      := make(chan Conn, 10)
    var sendMessage string


    go listener(listenConn, newConn_c)
    go dialer("129.241.187.147", ":6969", dialconn_c)

    go networking(newConn_c, generatedMsgs_c, receivedMsgs_c, dialconn_c)
    for {
        Scanf("%s", &sendMessage)
        generatedMsgs_c <- sendMessage+"EOL" 
        }
    }

