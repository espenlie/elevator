package networking

import (
    "misc"
    "net"
    "fmt"
    "time"
    "strings"
    "errors"
    "encoding/json"
    "elevator"
    "drivers"
)

type Con struct {
    Address *net.TCPConn
    Connect bool
}

var elevators = make(map[string]bool)
var orderlist = make([]Order, 0)
var statuslist = make(map[string]Status)
var connections =  make([]*net.TCPConn, 0)

func GetStatusList() map[string]Status {
    return statuslist
}

func GetOrderList() []Order {
    return orderlist
}

func PackNetworkMessage(message Networkmessage) []byte {
	send, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Could not pack message: ",err.Error())
	}
	return send
}

func UnpackNetworkMessage(pack []byte) Networkmessage{
	var message Networkmessage
    fmt.Println("UNPACKING: ", string(pack))
	err := json.Unmarshal(pack, &message)
	if err != nil {
		fmt.Println("Could not unpack message: ", err.Error())
	}
	return message
}

func InitUpdate(connection *net.TCPConn, myip string) {
    pack := make([]byte,1024)
    status := statuslist[myip]
    pack = PackNetworkMessage(Networkmessage{Order:Order{}, Status:status})
    time.Sleep(10*time.Millisecond)
    connection.Write(pack)
    for _,order := range orderlist {
            time.Sleep(10*time.Millisecond)
            pack = PackNetworkMessage(Networkmessage{Order:order,Status: Status{}})
            connection.Write(pack)
    }
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
                fmt.Println("before: ",connections)
//              connections = append(connections[:i], connections[i+1:]...)
                connections[len(connections)-1], connections[i], connections = nil, connections[len(connections)-1], connections[:len(connections)-1]
                fmt.Println("after: ",connections)
                return connections,nil
            }
        }
    return connections, errors.New("Connection not in slice") 
}

func Dialer2(connect_c chan Con, port string, elevators []misc.Elevator){
    local, _ := net.ResolveTCPAddr("tcp", "localhost"+port)
    localconn, _ := net.DialTCP("tcp",nil,local)
    connect_c <- Con{Address:localconn,Connect:true}
    fmt.Println("ELEV",elevators)
    for{

        cons := connections
        fmt.Println("CONS:",cons)
        elevatorloop:
	    for _,elevator := range elevators{
            fmt.Println("DIALER:",elevator)
            for _, connection := range cons {
                if strings.Split(connection.RemoteAddr().String(),":")[0] == elevator.Address {
                    continue elevatorloop
                }
            }
            fmt.Println("DIALING")
            raddr, err := net.ResolveTCPAddr("tcp",elevator.Address+port)
            dialConn, err := net.DialTCP("tcp", nil, raddr)
            if err != nil {
                fmt.Println("Dial ERROR: ", err)
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
            fmt.Println("Receiver:",err.Error())
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
    time.Sleep(10*time.Millisecond)
    for {
        _, err := connection.Write([]byte("KEEPALIVE"))
        if err != nil {
            error_c <- err.Error()
            return
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
                    go InitUpdate(connection.Address, myip)
                }else{
                    fmt.Println("Removing: ",connection)
                    remoteip := strings.Split(connection.Address.RemoteAddr().String(), ":")[0]
	                errorstate := Status{State: "ERROR", LastFloor: 0, Inhouse: false,Source: remoteip}
                    statuslist[remoteip] = errorstate
                    for i, con := range connections {
                        if con == connection.Address {
                            fmt.Println("before: ",connections)
            //              connections = append(connections[:i], connections[i+1:]...)
                            connections[len(connections)-1], connections[i], connections = nil, connections[len(connections)-1], connections[:len(connections)-1]
                        }
                    }
//                  _ , err := RemoveConnection(connections, connection.Address)
//                  if err != nil {
//                      error_c <- err.Error()
//                  }
                    connection.Address.Close()
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
                    fmt.Println("WILL wrtie: ",connection.RemoteAddr().String())
                    _, err := connection.Write(pack)
                    if err != nil{
                        error_c <- err.Error()
//                      connections_c <- Con{Address: connection, Connect: false}
                    }
                }
            }
            case err := <- error_c: {
                fmt.Println("ERROR: "+err)
            }
//          default:{
//              fmt.Println("My connections: ")
//              for _, connection := range connections {
//                  fmt.Println(connection.RemoteAddr().String())
//              }
//              time.Sleep(time.Microsecond)
//          }
        }
    }
}

func SendStatuslist(generatedMsgs_c chan Networkmessage) {
    myip := misc.GetLocalIP()
    mystatus := statuslist[myip]
    generatedMsgs_c <- Networkmessage{Order:Order{}, Status: mystatus}
}

func NewStatus(status Status, generatedMsgs_c chan Networkmessage) bool {
    for _, oldstat := range statuslist {
        if oldstat == status {
            return false
        }
    }
    generatedMsgs_c <- Networkmessage{Order: Order{}, Status: status}
    return true
}

func Neworder(generatedMsgs_c chan Networkmessage, order Order)bool{
    if order. Direction!=elevator.BUTTON_COMMAND{
        order.Source=""
    }
    for _, b := range orderlist {
        if b == order {
            return false
        }
    }
    generatedMsgs_c <- Networkmessage{Order:order, Status: Status{}}
    return true
}
