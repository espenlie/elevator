package networking

import (
    "net"
    "fmt"
    "time"
//  "strings"
    "errors"
    "io"
//  "os"

)
var elevators = make(map[string]bool)

type Con struct {
    Address *net.TCPConn
    Connect bool
}

func Orderdistr(generatedMsgs_c chan networking.Networkmessage){
    for{
        for floor, buttons := Button_channel_matrix{
            for button, channel:= buttons{
                if drivers.ReadBit(channel){
                    networking. Neworder(generatedMsgs_c, Order{Direction:button, Floor:floor+1, InOut:1, Source: myip})
                }
            }
        }
	time.Sleep(50 * time.Millisecond)
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

func Dialer2(connect_c chan Com, port string, dialconn_c chan *net.TCPConn){
	for{
		for elevator,status := range elevators{
			if !status {
                raddr, err := net.ResolveTCPAddr("tcp",elevator+port)
				dialConn, err := net.DialTCP("tcp", nil, raddr)
				if err != nil {
					fmt.Println(err)
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

func Listener2(conn *net.TCPListener, newConn_c chan *net.TCPConn, connect_c chan Com){
    for {
        newConn, err := conn.AcceptTCP()
        if err != nil {
            fmt.Println(err)
        }
        connect_c <- Com{Address:newConn, Connect:true}
        newConn_c <- newConn
    }
}

func Receiver2(conn *net.TCPConn, receivedMsgs_c chan Networkmessage){
    buf := make([]byte,1024)
    for {
        bit, err := conn.Read(buf[0:])
        if err != nil {
            fmt.Println(err)
            return
        }
        unpacked := UnpackNetworkMessage(buf,bit)
        receivedMsgs_c <- unpacked
    }
}
func IsAlive(connection *net.TCPConn, error_c chan string, connect_c chan Com) {
    for {
        connection.SetDeadline(time.Now().Add(3 * time.Microsecond))
        if _, err := connection.Write([]byte("a")); err != nil {
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
            error_c <- err.Error()
            return
        }
    }
}

func NetworkWrapper(elevators []Elevator) {
    var connections []*net.TCPConn
    listenaddr, _ := net.ResolveTCPAddr("tcp", "6969")
    lisenconn, _ := net.ListenTCP("tcp", listenaddr)
    connections_c := make(chan Con, 15)
    generatedmessages_c := make(chan Networkmessage, 10)
    receivedmessages_c := make(chan Networkmessage,15)
    error_c := make(chan string, 10)
    go Listener(listenconn, connections_c)
    go Orderdistr(generatedmessages_c)
    go Dialer(
    for {
        select {
            case connection := <- connections_c: {
                if connection.Connect {
                    connections = append(connections, newconnection)
                    go Receiver(connection)
//                  go IsAlive(newconnection, error_c, connect_c)
                }else{
                    connection.Close()
                    connections, err_= RemoveConnection(connections, connection)
                    if err != nil {
                        error_c <- err.Error()
                    }
                }

            }
            case received := <- receivedmessages_c: {
                if received.Order.Floor>0{
                    elevator.Elev_set_button_lamp(received.Order.Direction, received.Order.Floor, received.Order.receivedOut)
                    if received.Order.InOut==0{
                        received.Order.InOut=1
                        for i, b := range orderlist {
                            if b == received.Order {
                                orderlist = append(orderlist[:i], orderlist[i+1:]...)
                            }
                        }
                    }else{
                        orderlist=append(orderlist, received.Order)
                        fmt.Println(orderlist)
                    }
                }
            }
            case message := <- generatedmessages_c: {
                pack := make([]byte,1024)
                pack = PackNetworkMessage(message)
                for _,connection := range connections {
                    connection.SetDeadline(time.Now().Add(50 * time.Millisecond))
                    if _ err := connection.Write(packed); err != nil {
                        time.Sleep(100 * time.Millisecond)
                        error_c <- err.Error()
                        connections_c <- Con{Address: connection, Connect: false}
                        } 
            }
            case err := <- error_c: {
                fmt.Println("ERROR: "+err)
            }
            default: {
                time.Sleep(time.Second)
            }

        }
    }
} 

