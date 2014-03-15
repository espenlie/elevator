package networking

import (
    "elevator"
    "fmt"
    "net"
    "time"
    "drivers"
    "encoding/json"
    "misc"
)
var orderlist = make([]Order,0)
var statuslist = make(map[string]Status)

func PackNetworkMessage(message Networkmessage) []byte {
	send, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Could not pack message: ",err.Error())
	}
	return send
}

func UnpackNetworkMessage(pack []byte, bit int) Networkmessage{
	var message Networkmessage
	err := json.Unmarshal(pack[:bit], &message)
	if err != nil {
		fmt.Println("Could not unpack message: ", err.Error())
	}
	return message
}

func GenerateMessage(dir elevator.Elev_button, floor int, inout int, state string, lastfloor int, inhouse bool, source string) Networkmessage {
	s := Status{State: state, LastFloor: lastfloor, Inhouse: inhouse,Source:source}
	o := Order{Direction:dir, Floor:floor, InOut:inout}
	message := Networkmessage{Order:o,Status:s}
	return message
}

func SendStatuslist(generatedMsgs_c chan Networkmessage) {
    for _, status:= range statuslist {
        generatedMsgs_c <- GenerateMessage(elevator.BUTTON_CALL_UP,0,0,status.State, status.LastFloor,false,status.Source)
    }
}

func NewStatus(status Status, generatedMsgs_c chan Networkmessage) bool {
    for key, _ := range statuslist {
        if statuslist[key] == status {
            return false
        }
    }
    generatedMsgs_c <- GenerateMessage(elevator.BUTTON_CALL_UP,0,0,status.State, status.LastFloor,false,status.Source)
    return true
}


func neworder(generatedMsgs_c chan Networkmessage, order Order)bool{
    for _, b := range orderlist {
        if b == order {
            return false
        }
    }
    generatedMsgs_c <- GenerateMessage(order.Direction,order.Floor,order.InOut,"",-1,false,"")
    return true
}

func Orderdistr(generatedMsgs_c chan Networkmessage){
    for{
        if drivers.ReadBit(drivers.FLOOR_UP1){
            neworder(generatedMsgs_c, Order{Direction:elevator.BUTTON_CALL_UP, Floor:1, InOut:1})
        }
        if drivers.ReadBit(drivers.FLOOR_UP2){
            neworder(generatedMsgs_c, Order{Direction:elevator.BUTTON_CALL_UP, Floor:2, InOut:1})
        }
        if drivers.ReadBit(drivers.FLOOR_UP3){
            neworder(generatedMsgs_c, Order{Direction:elevator.BUTTON_CALL_UP, Floor:3, InOut:1})
        }
        if drivers.ReadBit(drivers.FLOOR_DOWN2){
            neworder(generatedMsgs_c, Order{Direction:elevator.BUTTON_CALL_DOWN, Floor:2, InOut:1})
        }
        if drivers.ReadBit(drivers.FLOOR_DOWN3){
            neworder(generatedMsgs_c, Order{Direction:elevator.BUTTON_CALL_DOWN, Floor:3, InOut:1})
        }
        if drivers.ReadBit(drivers.FLOOR_DOWN4){
            neworder(generatedMsgs_c, Order{Direction:elevator.BUTTON_CALL_DOWN, Floor:4, InOut:1})
        }
	time.Sleep(1 * time.Millisecond)
    }
}

//Receives messages from a connections and adds it to a channel
func Receiver(conn net.Conn, receivedMsgs_c chan Networkmessage){
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

//Listens for new connections. If any it accepts it and adds it to a connection channel
func Listener(conn *net.TCPListener, newConn_c chan net.Conn){
    for {
        newConn, err := conn.Accept()
        if err != nil {
            fmt.Println(err)
        }
        newConn_c <- newConn
    }
}


func Networking(newConn_c chan net.Conn, generatedMsgs_c chan Networkmessage, receivedMsgs_c chan Networkmessage, dialConn_c chan net.Conn) {
    var newConn net.Conn
//  var msg []Networkmessage
	connMap := make(map[string]net.Conn)
    for{
        select {
        case newConn = <- dialConn_c:
			connMap[newConn.LocalAddr().String()] = newConn
		case sendMsg := <- generatedMsgs_c:{
            packed := make([]byte,1024)
            packed = PackNetworkMessage(sendMsg)
			for _,connection := range connMap{
//              connection.SetLinger(100 * time.Millisecond)
				connection.Write(packed)
//              if err != nil {
//                  fmt.Println("KUK komputeren")
//              }
			}
		}
        case in := <-receivedMsgs_c:
            if in.Order.Floor>0{
                orderlist=append(orderlist, in.Order)
                elevator.Elev_set_button_lamp(in.Order.Direction, in.Order.Floor, in.Order.InOut) 
                fmt.Println(orderlist)
            }
            if in.Status.Source != "" {
                if in.Status.State == "INIT" {
                    SendStatuslist(generatedMsgs_c)
                }
                statuslist[in.Status.Source] = in.Status
                fmt.Println(statuslist)
            }
        case newConn := <- newConn_c:
            go Receiver(newConn, receivedMsgs_c)
        }

    }
}

//Dials all elevators in the map
func Dialer(elevators map[string]bool, port string, dialconn_c chan net.Conn){
    myip := misc.GetLocalIP()
    dialConn, err := net.Dial("tcp", "localhost"+port)
    if err != nil {
        fmt.Println(err)
    }
	dialconn_c <-dialConn
	for{
		for elevator,status := range elevators{
			if !status && elevator != myip{
				dialConn, err := net.Dial("tcp", elevator+port)
				if err != nil {
					fmt.Println(err)
				}else{
					elevators[elevator]=true
					dialconn_c <-dialConn
				}
			}
		}
    	time.Sleep(1000 * time.Millisecond)
	}
}
