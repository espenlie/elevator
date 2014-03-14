package networking

import (
    "elevator"
    "fmt"
    "net"
    "time"
    "drivers"
    "encoding/json"
    "bytes"
)
func PackNetworkMessage(message Networkmessage) []byte {
//  send := make([]byte,1024)
	send, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Could not pack message: ",err.Error())
	}
	return send
}

func UnpackNetworkMessage(pack []byte, bit int) Networkmessage{
	var message Networkmessage
    trimed := bytes.Trim(pack, "\x00")
	err := json.Unmarshal(trimed[:bit], &message)
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

func Orderdistr(generatedMsgs_c chan Networkmessage){
    for{
        if drivers.ReadBit(drivers.FLOOR_UP1){
            generatedMsgs_c <- GenerateMessage(elevator.BUTTON_CALL_UP,1,1,"",-1,false,"")
        }
        if drivers.ReadBit(drivers.FLOOR_UP2){
            generatedMsgs_c <- GenerateMessage(elevator.BUTTON_CALL_UP,2,1,"",-1,false,"")
        }
        if drivers.ReadBit(drivers.FLOOR_UP3){
            generatedMsgs_c <- GenerateMessage(elevator.BUTTON_CALL_UP,3,1,"",-1,false,"")
        }
        if drivers.ReadBit(drivers.FLOOR_DOWN2){
            generatedMsgs_c <- GenerateMessage(elevator.BUTTON_CALL_DOWN,2,1,"",-1,false,"")
        }
        if drivers.ReadBit(drivers.FLOOR_DOWN3){
            generatedMsgs_c <- GenerateMessage(elevator.BUTTON_CALL_DOWN,3,1,"",-1,false,"")
        }
        if drivers.ReadBit(drivers.FLOOR_DOWN4){
            generatedMsgs_c <- GenerateMessage(elevator.BUTTON_CALL_DOWN,4,1,"",-1,false,"")
        }
	time.Sleep(5 * time.Millisecond)
    }
}

//Receives messages from a connections and adds it to a channel
func Receiver(conn net.Conn, receivedMsgs_c chan Networkmessage){
    buf := make([]byte,1024)
//  var buf []byte
//  conn.SetReadBuffer(1024)
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
//          fmt.Println(sendMsg)
            packed = PackNetworkMessage(sendMsg)
			for _,connection := range connMap{
				connection.Write(packed)
//              if err != nil {
//                  fmt.Println("Error writing: ", err.Error())
//              }
//              fmt.Println(mess)
			}
		}
        case in := <-receivedMsgs_c:
            fmt.Println(in)
            if in.Order.Floor>0{
//              orderlist=append(orderlist, in.Order)
                elevator.Elev_set_button_lamp(in.Order.Direction, in.Order.Floor, in.Order.InOut) 

            }
        case newConn := <- newConn_c:
            go Receiver(newConn, receivedMsgs_c)
        }

    }
}

//Dials all elevators in the map
func Dialer(elevators map[string]bool, port string, dialconn_c chan net.Conn){
	for{
		for elevator,status := range elevators{
//			Println(elevators)
			if !status{
				dialConn, err := net.Dial("tcp", elevator+port)
				if err != nil {
					fmt.Println(err)
				}else{
					elevators[elevator]=true
					dialconn_c <-dialConn
				}
			}
		}
    	time.Sleep(10000 * time.Millisecond)
	}
}
