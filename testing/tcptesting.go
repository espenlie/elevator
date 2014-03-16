package main

import (
    "net"
    "fmt"
    "time"
    "strings"
//  "io"

)

var elevator = make(map[string]bool)

type Com struct {
    Address string
    Connect bool
}


func main() {
    elevator["129.241.187.156"]=false
    elevator["129.241.187.161"]=false
    elevator["129.241.187.158"]=false
//	connectionmap := make(map[string]*net.TCPConn)
    var connections []*net.TCPConn
    connections_c := make(chan *net.TCPConn, 10)
    message_c     := make(chan []byte, 10)
    error_c       := make(chan string, 10)
    connect_c     := make(chan Com,10)
    listenaddr, _ := net.ResolveTCPAddr("tcp", ":5555")
    listenconn, _ := net.ListenTCP("tcp",listenaddr)
    go Listener(listenconn, connections_c, connect_c)
    go IsAlive(connections, error_c, connect_c)
    go Dialer(connect_c,":5555", connections_c)
    for {
        select {
            case newconnection := <- connections_c  :{
                fmt.Println("New connection",newconnection.LocalAddr().String())
                connections = append(connections, newconnection)
//              connectionmap[newconnection.LocalAddr().String()] = newconnection
            }
            case newmessage := <- message_c :{
                fmt.Println(string(newmessage))
            }
            case errorm := <- error_c :{
                fmt.Println("Errormessage:"+errorm)
            }
            case lost := <- connect_c :{
                elevator[lost.Address]=lost.Connect
            }

            default :{
                time.Sleep(1*time.Second)
                fmt.Println(connections)
            }
        }
    }

}

func IsAlive(connections []*net.TCPConn, error_c chan string, connect_c chan Com) {
    for{
        for i, connection := range connections {
//          var test net.PacketConn
//          var buff []byte
//          reads, err := connection.Read(buff)
//          err := connection.SetKeepAlive(true)
//          connection.SetWriteDeadline(time.Now().Add(time.Second))
//          _, err := test.WriteTo([]byte("test"),connection.RemoteAddr())
            fmt.Println(connection.RemoteAddr().String())
            connection.SetDeadline(time.Now().Add(time.Second))
            _, err := connection.Write([]byte("test"))
//          err := connection.SetLinger(1)
//          p, _ := connection.Write([]byte("test"))
//          connection.SetKeepAlive(true)
//          err := connection.SetKeepAlivePeriod(time.Second)
            if err != nil {
//          if reads,err := connection.Read(buff); err == io.EOF {
                connection.Close() 
                connect_c <- Com{Address:connection.RemoteAddr().String(),Connect:false}
//              elevator[i]=false
//              fmt.Printf("%v%v",j,p)
                error_c <- err.Error()
                connections[i] = connections[len(connections)-1]
//              delete(connections,i)
//          }else{
//              connection.SetReadDeadline(time.Time{})
//              fmt.Println(reads)
            }
        }
        time.Sleep(1*time.Second)
    }
}
func Dialer(connect_c chan Com, port string, dialconn_c chan *net.TCPConn){
	for{
		for elevator,status := range elevator{
			if !status {
                raddr, err := net.ResolveTCPAddr("tcp",elevator+port)
				dialConn, err := net.DialTCP("tcp", nil, raddr)
				if err != nil {
//					fmt.Println(err)
				}else{
                    connect_c <- Com{Address:elevator,Connect:true}
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
        connect_c <- Com{Address:strings.Split(newConn.RemoteAddr().String(),":")[0], Connect:true}
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
}
*/
