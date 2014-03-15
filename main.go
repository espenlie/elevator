package main
import (
    ."fmt"
    ."net"
//  ."strings"
//	"strconv"
    "time"
    "drivers"
	"misc"
	"networking"
	"elevator"
)
func nextorder(myip string)int{
	statuslist := networking.GetStatusList()
	orderlist := networking.GetOrderList()

	for _,order := range orderlist{
		for _,elev := range statuslist{
			if elev.State=="IDLE" && elev.LastFloor==order.Floor && elev.Source==myip{
				return order.Floor
			}
		}
		for _,elev := range statuslist{
			if ((elev.State=="UP" && elev.LastFloor==order.Floor-1) || (elev.State=="DOWN" && elev.LastFloor==order.Floor+1))&& elev.Source==myip{
				return order.Floor
			}
		}
		for _,elev := range statuslist{
			if elev.State=="IDLE" && (elev.LastFloor==order.Floor+1 || elev.LastFloor==order.Floor-1) && elev.Source==myip{
				return order.Floor
			}
		}
		for _,elev := range statuslist{
			if ((elev.State=="UP" && elev.LastFloor==order.Floor-2) || (elev.State=="DOWN" && elev.LastFloor==order.Floor+2))&& elev.Source==myip{
				return order.Floor
			}
		}
		for _,elev := range statuslist{
			if elev.State=="IDLE" && (elev.LastFloor==order.Floor+2 || elev.LastFloor==order.Floor-2) && elev.Source==myip{
				return order.Floor
			}
		}
		for _,elev := range statuslist{
			if ((elev.State=="UP" && elev.LastFloor==order.Floor-3) || (elev.State=="DOWN" && elev.LastFloor==order.Floor+3))&& elev.Source==myip{
				return order.Floor
			}
		}
		for _,elev := range statuslist{
			if elev.State=="IDLE" && (elev.LastFloor==order.Floor+3 || elev.LastFloor==order.Floor-3) && elev.Source==myip{
				return order.Floor
			}
		}

	}
	return 0
}

func main() {

	myip := misc.GetLocalIP()
	Println(myip)
	go elevator.FloorUpdater()

//	var conf misc.Config
	conf := misc.LoadConfig("/home/student/LL/elevator/config/conf.json")

    connections         := make(map[string]bool)

	for _,elevator :=range conf.Elevators{
//		Println(elevator.Address)
		connections[elevator.Address]=false
	}

    listenAddr, _ := ResolveTCPAddr("tcp", ":6969")
    listenConn, _ := ListenTCP("tcp", listenAddr)
    receivedMsgs_c  := make(chan networking.Networkmessage)
    generatedMsgs_c  := make(chan networking.Networkmessage)
    newConn_c       := make(chan Conn, 10)
    dialConn_c      := make(chan Conn, 10)

    go networking.Networking(newConn_c, generatedMsgs_c, receivedMsgs_c, dialConn_c)
//	statuslist[myip]=networking.Status{State:"UP",LastFloor:1,Source:myip}
//	takeorder(orderlist, statuslist, myip)
	go networking.Listener(listenConn, newConn_c)
	go networking.Dialer(connections, conf.Default_Dial_Port, dialConn_c)
	go networking.Orderdistr(generatedMsgs_c)
//	for {
//      Scanf("%s", &sendMessage)
//      generatedMsgs_c <- sendMessage+"EOL" 

//      }

	state := "INIT"
//	var floor int
	var mystatus networking.Status
	mystatus.Source=myip
	mystatus.State=state
	mystatus.LastFloor=elevator.Current_floor()
	time.Sleep(1500 * time.Millisecond)
	for{
//		Println("State: ", state)
		switch state {
			case "INIT":{
				drivers.IoInit()
				elevator.Elev_init()
//				Println(mystatus)
				time.Sleep(10 * time.Millisecond)
				networking.NewStatus(mystatus, generatedMsgs_c)
				elevator.Elev_set_speed(-300)
				if elevator.Elev_get_floor_sensor_signal()!=-1{
					elevator.Elev_set_speed(0)
					state="IDLE"
				}
				fallthrough;
				}
			case "IDLE":{
				Println(nextorder(myip))
				time.Sleep(10 * time.Millisecond)
//				status.State = state
//				status.LastFloor = floor
				fallthrough;
			}
			case "UP":{
				fallthrough;
			}
			case "DOWN":{
				fallthrough;
			}
			case "DOOR_OPEN":{
				fallthrough;
			}
			case "ERROR":{
				fallthrough;
			}
			default:{
				mystatus.State=state
				mystatus.LastFloor=elevator.Current_floor()
				networking.NewStatus(mystatus, generatedMsgs_c)
//				generatedMsgs_c <- networking.GenerateMessage(elevator.BUTTON_CALL_UP,0,0,mystatus.State, mystatus.LastFloor,false,mystatus.Source)
			}
		}
	}
}
