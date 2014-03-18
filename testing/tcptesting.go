package main

import (
    "net"
    "fmt"
    "time"
    "strings"
    "errors"
    "io"
//  "os"

)

var elevator = make(map[string]bool)

type Com struct {
    Address *net.TCPConn
    Connect bool
}

func main() {
//  elevator["193.35.52.151"]=false
//  elevator["193.35.52.194"]=false
//  elevator["193.35.52.234"]=false
    elevator["129.241.187.153"]=false
    elevator["129.241.187.140"]=false
    elevator["129.241.187.161"]=false
    var connections []*net.TCPConn
    connections_c := make(chan *net.TCPConn, 10)
    message_c     := make(chan []byte, 10)
    error_c       := make(chan string, 10)
    connect_c     := make(chan Com,10)

    listenaddr, _ := net.ResolveTCPAddr("tcp", ":5555")
    listenconn, _ := net.ListenTCP("tcp",listenaddr)
    go Listener(listenconn, connections_c, connect_c)
    go Dialer(connect_c,":5555", connections_c)
    for {
        select {
            case newconnection := <- connections_c  :{
                fmt.Println("New connection",newconnection.LocalAddr().String())
//              newconnection.SetDeadline(time.Now().Add(1*time.Second))
                connections = append(connections, newconnection)
                go IsAlive(newconnection, error_c, connect_c)
            }
            case newmessage := <- message_c :{
                fmt.Println(string(newmessage))
            }
            case errorm := <- error_c :{
                fmt.Println("Errormessage:"+errorm)
            }
            case lost := <- connect_c :{
                elevator[strings.Split(lost.Address.RemoteAddr().String(),":")[0]]=lost.Connect
                if !lost.Connect{
                    connections, _ = RemoveConnection(connections, lost.Address)
//                  if err != nil {
//                      fmt.Println(err)
//                  }
//                  fmt.Println(connections)
                }
            }

            default :{
                time.Sleep(1*time.Second)
                fmt.Println(connections)
            }
        }
    }
}
//func isTimeout(err error) bool {
//  e, ok := err.(Error)
//  return ok && e.Timeout()
//}

func RemoveConnection(connections []*net.TCPConn, connection *net.TCPConn) ([]*net.TCPConn, error) {
        for i, con := range connections {
            if con == connection {
                connections = append(connections[:i], connections[i+1:]...)
                return connections,nil
            }
        }
    return connections, errors.New("Connection not in slice") 
}


func IsAlive(connection *net.TCPConn, error_c chan string, connect_c chan Com) {
    for{
        connection.SetDeadline(time.Now().Add(3 * time.Microsecond))
//      connection.Write([]byte("test"))
//      var buf []byte
//      if _, err := connection.Read(buf[:]); err != nil {
        
        if _, err := connection.Write([]byte("a")); err != nil {
//          if isTimout(err){
//              fmt.Println("TIMEOUT")
//          }
//          fmt.Println(err.Error())
            time.Sleep(time.Second)
            if opErr, ok := err.(*net.OpError); ok{
                if opErr.Timeout() {
                    fmt.Println("TIMEOUT")
                }
                if opErr.Temporary() {
                    fmt.Println("TEMPORARY")
                }
            }
//          if err == io.EINVAL {
//              fmt.Println("EINVAL")
//          }
            if err == io.EOF {
                fmt.Println("EOF")
            }
            connection.Close()
            connect_c <- Com{Address:connection,Connect:false}
//          error_c <- err.Error()
            return
        }
//      if err != nil {
//          connection.Close()
//          connect_c <- Com{Address:connection,Connect:false}
//          error_c <- err.Error()
//          break
//      }
        time.Sleep(1*time.Second)
    }
}
func Dialer(connect_c chan Com, port string, dialconn_c chan *net.TCPConn){
	for{
		for elevator,status := range elevator{
			if !status {
                raddr, err := net.ResolveTCPAddr("tcp",elevator+port)
                fmt.Println("Dialing: "+raddr.String())
				dialConn, err := net.DialTCP("tcp", nil, raddr)
				if err != nil {
//					fmt.Println(err)
				}else{
                    connect_c <- Com{Address:dialConn,Connect:true}
                    fmt.Println("Adding: ",dialConn)
					dialconn_c <-dialConn
				}
			}
		}
	    time.Sleep(1000 * time.Millisecond)
	}
}

func Listener(conn *net.TCPListener, newConn_c chan *net.TCPConn, connect_c chan Com){
    for {
        newConn, err := conn.AcceptTCP()
        if err != nil {
            fmt.Println(err)
        }
        fmt.Println("New shit")
        connect_c <- Com{Address:newConn, Connect:true}
        newConn_c <- newConn
    }
}
/*
func Receiver(conn *net.TCPConn, receivedMsgs_c chan Networkmessage){
    buf := make([]byte,1024)
    for {
        bit, err := conn.Read(buf[0:])
        if err != nil {
            fmt.Println(err)
            return
        }
        receivedMsgs_c <- bit
    }
}*/

