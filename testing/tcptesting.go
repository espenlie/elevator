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
    elevator["129.241.187.141"]=false
    elevator["129.241.187.143"]=false
    elevator["129.241.187.150"]=false
//  elevator["129.241.187.161"]=false
    var connections []*net.TCPConn
    connections_c := make(chan *net.TCPConn, 10)
    message_c     := make(chan []byte, 10)
    receive_c     := make(chan string, 10)
    error_c       := make(chan string, 10)
    connect_c     := make(chan Com,10)

    listenaddr, _ := net.ResolveTCPAddr("tcp", ":6666")
    listenconn, _ := net.ListenTCP("tcp",listenaddr)
    go Listener(listenconn, connections_c, connect_c)
    go Dialer(connect_c,":6666", connections_c)
    for {
        select {
            case newconnection := <- connections_c  :{
                fmt.Println("New connection",newconnection.LocalAddr().String())
//              newconnection.SetDeadline(time.Now().Add(1*time.Second))
                connections = append(connections, newconnection)
//              newconnection.SetKeepAlive(true)
//              newconnection.SetKeepAlivePeriod(2000*time.Millisecond)
                go Receiver(newconnection, receive_c, connect_c)
//              go IsAlive(newconnection, error_c, connect_c)
            }
            case newmessage := <-message_c :{
                fmt.Println(string(newmessage))
            }
            case errorm := <- error_c :{
                fmt.Println("Errormessage:"+errorm)
            }
            case in := <- receive_c :{
                fmt.Println("INCOMING: ",in)
            }
            case lost := <- connect_c :{
                index := strings.Split(lost.Address.RemoteAddr().String(),":")[0]
                elevator[index]=lost.Connect
                if !lost.Connect{
                    connections, _ = RemoveConnection(connections, lost.Address)
                }
            }

            default :{
                time.Sleep(10*time.Millisecond)
                fmt.Println(connections)
            }
        }
    }
}

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
        connection.SetWriteDeadline(time.Now().Add(30 * time.Microsecond))
        connection.SetKeepAlive(true)
        connection.SetKeepAlivePeriod(10*time.Millisecond)
        if _, err := connection.Write([]byte("hei")); err != nil {
            time.Sleep(time.Second)
            if opErr, ok := err.(*net.OpError); ok{
                if opErr.Timeout() {
                    fmt.Println("TIMEOUT")
                }
                if opErr.Temporary() {
                    fmt.Println("TEMPORARY")
                }
            }
            if err == io.EOF {
                fmt.Println("EOF")
            }
            connection.Close()
            connect_c <- Com{Address:connection,Connect:false}
            return
        }
        time.Sleep(500*time.Millisecond)
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
					fmt.Println(err)
				}else{
                    connect_c <- Com{Address:dialConn,Connect:true}
                    fmt.Println("Adding: ",dialConn)
                    fmt.Println(dialConn.RemoteAddr().String())
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
            fmt.Println("Listen error: ",err)
        }
        connect_c <- Com{Address:newConn, Connect:true}
        newConn_c <- newConn
    }
}

func Receiver(conn *net.TCPConn, receivedMsgs_c chan string, connect_c chan Com){
    buf := make([]byte,1024)
    for {
        bit, err := conn.Read(buf[0:])
        conn.SetReadDeadline(time.Now().Add(5*time.Nanosecond))
        if err != nil {
            fmt.Println("READ ERR: ",err)
            conn.Close()
            connect_c <- Com{Address:conn,Connect:false}
            return
        }
        receivedMsgs_c <- string(bit)
    }
}

