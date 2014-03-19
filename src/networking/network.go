package networking

import (
    "misc"
//  "bytes"
    "net"
    "fmt"
    "time"
    "strings"
    "errors"
    "encoding/json"
    "elevator"
//  "os"
    "drivers"

)
var elevators = make(map[string]bool)
var orderlist = make([]Order, 0)
//var insidelist = make([]Order, 0)
var statuslist = make(map[string]Status)
var connections =  make([]*net.TCPConn, 0)

func GetStatusList() map[string]Status {
    return statuslist
}
//func GetInsideList() []Order {
//  return insidelist
//}
func GetOrderList() []Order {
    return orderlist
}

func PackNetworkMessage(message Networkmessage) []byte {
//  fmt.Println("PACKING: ", message)
	send, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Could not pack message: ",err.Error())
	}
	return send
}

func UnpackNetworkMessage(pack []byte) Networkmessage{
	var message Networkmessage
//  fmt.Println("UNPACKING: ", string(pack))
	err := json.Unmarshal(pack, &message)
	if err != nil {
		fmt.Println("Could not unpack message: ", err.Error())
	}
	return message
}

type Con struct {
    Address *net.TCPConn
    Connect bool
}

func Orderdistr(generatedMsgs_c chan Networkmessage, myip string){
    var butt elevator.Elev_button
    for{
        for floor, buttons := range elevator.Button_channel_matrix{
            for button, channel:= range buttons{
                if drivers.ReadBit(channel){
                    if button == 0 {
                        butt = elevator.BUTTON_CALL_UP
                    }else if button == 1 {
                        butt = elevator.BUTTON_CALL_DOWN
                    }else{
                        butt = elevator.BUTTON_COMMAND
                    }
                    Neworder(generatedMsgs_c, Order{Direction:butt, Floor:floor+1, InOut:1, Source: myip})
                    time.Sleep(time.Millisecond)
                }
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

func Dialer2(connect_c chan Con, port string, elevators []misc.Elevator){
    for{
        elevatorloop:
	    for _,elevator := range elevators{
            for _, connection := range connections {
//              fmt.Println(connection)
//              fmt.Println(connection.RemoteAddr().String())
                if strings.Split(connection.RemoteAddr().String(),":")[0] == elevator.Address {
                    continue elevatorloop
                }
            }
            raddr, err := net.ResolveTCPAddr("tcp",elevator.Address+port)
            dialConn, err := net.DialTCP("tcp", nil, raddr)
            if err != nil {
//              fmt.Println("Dial ERROR: ", err)
            }else{
                connect_c <- Con{Address:dialConn, Connect:true}
                fmt.Println("Adding: ",dialConn.RemoteAddr().String())
            }
		}
	    time.Sleep(1000 * time.Millisecond)
	}
}

func Listener2(conn *net.TCPListener, connect_c chan Con){
    for {
        newConn, err := conn.AcceptTCP()
        if err != nil {
            fmt.Println("AcceptERROR: ", err)
        }
        connect_c <- Con{Address:newConn, Connect:true}
    }
}

func Receiver2(conn *net.TCPConn, receivedMsgs_c chan Networkmessage, connections_c chan Con){
    buf := make([]byte,1024)
    keepalivebyte := []byte("KEEPALIVE")
    receiverloop:
    for {
        conn.SetReadDeadline(time.Now().Add(2*time.Second))
        bit, err := conn.Read(buf[0:])
        if err != nil {
            fmt.Println(err.Error())
            connections_c <- Con{Address:conn,Connect:false}
            return
        }
        if string(buf[:bit]) == string(keepalivebyte){
            continue receiverloop
        }
        unpacked := UnpackNetworkMessage(buf[:bit])
        receivedMsgs_c <- unpacked
    }
}

func SendAliveMessages(connection *net.TCPConn, error_c chan string) {
    for {
        _, err := connection.Write([]byte("KEEPALIVE"))
        if err != nil {
            error_c <- err.Error()
        }
        time.Sleep(time.Second)
    }
}

func NetworkWrapper(conf misc.Config, myip string, generatedmessages_c chan Networkmessage) {
    listenaddr, _ := net.ResolveTCPAddr("tcp", ":5555")
    listenconn, _ := net.ListenTCP("tcp", listenaddr)
    connections_c := make(chan Con, 15)
    receivedmessages_c := make(chan Networkmessage,15)
    error_c := make(chan string, 10)
    go Listener2(listenconn, connections_c)
    go Orderdistr(generatedmessages_c, myip)
    go Dialer2(connections_c, ":5555", conf.Elevators)
    for {
        select {
            case connection := <- connections_c: {
                if connection.Connect {
                    connections = append(connections, connection.Address)
                    go Receiver2(connection.Address, receivedmessages_c, connections_c)
                    go SendAliveMessages(connection.Address, error_c)
                }else{
                    fmt.Println("Removing: ",connection)
                    connection.Address.Close()
                    _ , err := RemoveConnection(connections, connection.Address)
                    if err != nil {
                        error_c <- err.Error()
                    }
                }

            }
            case received := <- receivedmessages_c: {

                if received.Order.Floor>0{
                    if !((received.Order.Direction == elevator.BUTTON_COMMAND) && (received.Order.Source != myip)) {
                        elevator.Elev_set_button_lamp(received.Order.Direction, received.Order.Floor, received.Order.InOut)
                    }
                    if received.Order. Direction!=elevator.BUTTON_COMMAND{
                        received.Order.Source=""
                    }
                    if received.Order.InOut==0{
                        received.Order.InOut=1
                            for i, b := range orderlist{
                                if b == received.Order { 
                                    orderlist = append(orderlist[:i], orderlist[i+1:]...)
                                }
                            }
                    }else{
                        AddedBefore:=false
                        for _, b := range orderlist {
                            if b == received.Order {
                                AddedBefore = true
                            }                            }
                        if !AddedBefore{
                            orderlist=append(orderlist, received.Order)
                        }
                    }
                }
                if received.Status.Source != "" {
                    statuslist[received.Status.Source] = received.Status
                }
            }
            case message := <- generatedmessages_c: {
                pack := make([]byte,1024)
                pack = PackNetworkMessage(message)
                for _,connection := range connections {
//                  connection.SetDeadline(time.Now().Add(50 * time.Millisecond))
                    connection.Write(pack)
//                  time.Sleep(100 * time.Millisecond)
//                  if err != nil{    
//                      error_c <- err.Error()
//                      connections_c <- Con{Address: connection, Connect: false}
//                  } 
                }
            }
            case err := <- error_c: {
                fmt.Println("ERROR: "+err)
            }
            default: {
//              time.Sleep(time.Second)
            }

        }
    }
} 
func GenerateMessage(dir elevator.Elev_button, floor int, inout int, state string, lastfloor int, inhouse bool, source string) Networkmessage {
	s := Status{State: state, LastFloor: lastfloor, Inhouse: inhouse,Source:source}
	o := Order{Direction:dir, Floor:floor, InOut:inout}
	message := Networkmessage{Order:o,Status:s}
	return message
}

func SendStatuslist(generatedMsgs_c chan Networkmessage) {
    myip := misc.GetLocalIP()
    mystatus := statuslist[myip]
    generatedMsgs_c <- GenerateMessage(elevator.BUTTON_CALL_UP,0,0,mystatus.State, mystatus.LastFloor,false,mystatus.Source)
}

func NewStatus(status Status, generatedMsgs_c chan Networkmessage) bool {
//  fmt.Println("statuslist: ", statuslist)
//  fmt.Println("Newstatus: ", status)
    for _, oldstat := range statuslist {
        if oldstat == status {
            return false
        }
    }
//  fmt.Println("Sending status")
    generatedMsgs_c <- GenerateMessage(elevator.BUTTON_CALL_UP,0,0,status.State, status.LastFloor,false,status.Source)
    return true
}


func Neworder(generatedMsgs_c chan Networkmessage, order Order)bool{
//  fmt.Println("Orderlist: ",orderlist)
//  fmt.Println("Order: ", order)
    if order. Direction!=elevator.BUTTON_COMMAND{
        order.Source=""
    }
    for _, b := range orderlist {
        if b == order {
            return false
        }
    }
//  fmt.Println("Sending order")
    generatedMsgs_c <- Networkmessage{Order:order, Status: Status{"",-1,false,""}}
    return true
}
