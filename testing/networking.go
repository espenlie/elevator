package networking

import (
    "elevator"
    "fmt"
    "net"
    "time"
//  "drivers"
    "encoding/json"
    "misc"
)



//func Orderdistr(generatedMsgs_c chan Networkmessage){
//  for{
//      if drivers.ReadBit(drivers.FLOOR_UP1){
//          Neworder(generatedMsgs_c, Order{Direction:elevator.BUTTON_CALL_UP, Floor:1, InOut:1})
//      }
//      if drivers.ReadBit(drivers.FLOOR_UP2){
//          Neworder(generatedMsgs_c, Order{Direction:elevator.BUTTON_CALL_UP, Floor:2, InOut:1})
//      }
//      if drivers.ReadBit(drivers.FLOOR_UP3){
//          Neworder(generatedMsgs_c, Order{Direction:elevator.BUTTON_CALL_UP, Floor:3, InOut:1})
//      }
//      if drivers.ReadBit(drivers.FLOOR_DOWN2){
//          Neworder(generatedMsgs_c, Order{Direction:elevator.BUTTON_CALL_DOWN, Floor:2, InOut:1})
//      }
//      if drivers.ReadBit(drivers.FLOOR_DOWN3){
//          Neworder(generatedMsgs_c, Order{Direction:elevator.BUTTON_CALL_DOWN, Floor:3, InOut:1})
//      }
//      if drivers.ReadBit(drivers.FLOOR_DOWN4){
//          Neworder(generatedMsgs_c, Order{Direction:elevator.BUTTON_CALL_DOWN, Floor:4, InOut:1})
//      }
//	time.Sleep(50 * time.Millisecond)
//  }
//}

//Receives messages from a connections and adds it to a channel
func Receiver(conn *net.TCPConn, receivedMsgs_c chan Networkmessage){
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
func Listener(conn *net.TCPListener, newConn_c chan *net.TCPConn){
    for {
        newConn, err := conn.AcceptTCP()
        if err != nil {
            fmt.Println(err)
        }
        newConn_c <- newConn
    }
}


func Networking(newConn_c chan *net.TCPConn, generatedMsgs_c chan Networkmessage, receivedMsgs_c chan Networkmessage, dialConn_c chan *net.TCPConn) {
//  var msg []Networkmessage
	connMap := make(map[string]*net.TCPConn)
    for {
        select {
        case newConn := <- dialConn_c:
            newConn.SetKeepAlive(true)
            err := newConn.SetKeepAlivePeriod(100 * time.Millisecond)
            time.Sleep(time.Second)
            if err != nil {
                fmt.Println(err)
                newConn.Close()
            }
            newConn.SetLinger(1)
			connMap[newConn.LocalAddr().String()] = newConn
		case sendMsg := <- generatedMsgs_c:{
//          fmt.Println("Sendt: ", sendMsg)
            packed := make([]byte,1024)
            packed = PackNetworkMessage(sendMsg)
			for _,connection := range connMap{
//              err := connection.SetLinger(1)
                err := connection.SetWriteDeadline(time.Now().Add(1000 * time.Millisecond))
				connection.Write(packed)
                if err != nil {
                    fmt.Println("KUK komputeren")
                }
			}
		}
        case in := <-receivedMsgs_c:
//          fmt.Println("Motatt: ", in)
            if in.Order.Floor>0{
                elevator.Elev_set_button_lamp(in.Order.Direction, in.Order.Floor, in.Order.InOut)
                if in.Order.InOut==0{
                    in.Order.InOut=1
                    for i, b := range orderlist {
                        if b == in.Order {
                            orderlist = append(orderlist[:i], orderlist[i+1:]...)
                        }
                    }
                }else{

                    orderlist=append(orderlist, in.Order)
                    fmt.Println(orderlist)
                }
            }
            if in.Status.Source != "" {
                fmt.Println("Statuslist before updating: ", statuslist)
                fmt.Println("Adding: ", in.Status)
                statuslist[in.Status.Source] = in.Status
                fmt.Println("Statuslist after updating: ", statuslist)
                if in.Status.State == "INIT" && in.Status.Source != misc.GetLocalIP(){
                    go SendStatuslist(generatedMsgs_c)
                }
            }
        case newConn := <- newConn_c:
            go Receiver(newConn, receivedMsgs_c)
        }

    }
}

//Dials all elevators in the map
func Dialer(elevators map[string]bool, port string, dialconn_c chan *net.TCPConn){
    myip := misc.GetLocalIP()
//  local := net.TCPAddr{IP:[]byte("127.0.0.1"),Port:port}
    local, err := net.ResolveTCPAddr("tcp","localhost"+port)
    if err != nil {
        fmt.Println(err)
    }
    local2, err := net.ResolveTCPAddr("tcp","localhost"+":1337")
    if err != nil {
        fmt.Println(err)
    }
    dialConn, err := net.DialTCP("tcp4",local2, local)
    elevators[myip]=true
    if err != nil {
        fmt.Println(err)
    }
	dialconn_c <-dialConn
	for{
		for elevator,status := range elevators{
			if !status {
                raddr, err := net.ResolveTCPAddr("tcp",elevator+port)
				dialConn, err := net.DialTCP("tcp4", nil, raddr)
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
