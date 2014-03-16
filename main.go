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


func nextorder(myip string, connections map[string]bool)networking.Order{
	statuslist := networking.GetStatusList()
	orderlist := networking.GetOrderList()
	orderloop:
	for _,order := range orderlist{
		for i := 0; i < elevator.N_FLOORS; i++ {
			for elevator,_ :=range connections{
				if val,ok := statuslist[elevator]; ok{
					Println(val)
					if statuslist[elevator].State=="IDLE" && statuslist[elevator].LastFloor==order.Floor+i && statuslist[elevator].Inhouse==false{
						if statuslist[elevator].Source==myip{
							return order
						}else{
							continue orderloop
						}
					}
				if i!=0{
						if ((statuslist[elevator].State=="UP" && statuslist[elevator].LastFloor==order.Floor-i) || (statuslist[elevator].State=="DOWN" && statuslist[elevator].LastFloor==order.Floor+i)){
							if statuslist[elevator].Source==myip{
								return order
							}else{
							continue orderloop
							}
						}
					}
				}
			}
		}
	}
	return networking.Order{Direction:0,Floor:0,InOut:0}
}

func nextstate(myip string, connections map[string]bool)string{
	next := nextorder(myip, connections)
	if next.Floor>elevator.Current_floor(){
		return "UP"
	}else if (next.Floor<elevator.Current_floor() && next.Floor!=0){
		return "DOWN"
	}else if next.Floor==elevator.Current_floor(){
		return "DOOR_OPEN"
	}else{
		return "IDLE"
	}
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
    receivedMsgs_c  := make(chan networking.Networkmessage, 10)
    generatedMsgs_c  := make(chan networking.Networkmessage, 10)
    newConn_c       := make(chan *TCPConn, 10)
    dialConn_c      := make(chan *TCPConn, 10)

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
	var taken networking.Order
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
				networking.NewStatus(mystatus, generatedMsgs_c)
				elevator.Elev_set_speed(-300)
				if elevator.Elev_get_floor_sensor_signal()!=-1{
					elevator.Elev_set_speed(0)
					state="IDLE"
				}
			}
			case "IDLE":{
				elevator.Elev_set_speed(0)
//				Println(nextorder(myip))
				state =nextstate(myip, connections)
			}
			case "UP":{
				elevator.Elev_set_speed(300)
				state =nextstate(myip, connections)
			}
			case "DOWN":{
				elevator.Elev_set_speed(-300)
				state =nextstate(myip, connections)
			}
			case "DOOR_OPEN":{
				elevator.Elev_set_door_open_lamp(1)
				taken.InOut=0
				time.Sleep(50 * time.Millisecond)
				networking.Neworder(generatedMsgs_c, taken)
				elevator.Elev_set_speed(0)
				time.Sleep(3000 * time.Millisecond)
				elevator.Elev_set_door_open_lamp(0)
				state="IDLE"
			}
			case "ERROR":{
			}
		}
//		statuslist := networking.GetStatusList()
//		orderlist := networking.GetOrderList()
//		Println(statuslist)
//		Println(orderlist)
//		Println(state)
//		Println(nextorder(myip))
		time.Sleep(50 * time.Millisecond)
		Println(state)
		mystatus.State=state
		mystatus.LastFloor=elevator.Current_floor()
		networking.NewStatus(mystatus, generatedMsgs_c)
//		generatedMsgs_c <- networking.GenerateMessage(elevator.BUTTON_CALL_UP,0,0,mystatus.State, mystatus.LastFloor,false,mystatus.Source)
	}
}
