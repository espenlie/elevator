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
	statelist := networking.GetStatusList()
	ordlist := networking.GetOrderList()
//	Println("statuslist: ", statuslist)
//	Println("orderlist: ", orderlist)
//	Println("Connections: ", connections)
//	Println("MYIP: ", myip)
	orderloop:
	for _,order := range ordlist{
		for i := 0; i < elevator.N_FLOORS; i++ {
			for elevator,_ :=range connections{
				if status,ok := statelist[elevator]; ok{
					Println("Checking:  ", statelist[elevator], "take order: ", order)
					if ((status.State=="UP" && status.LastFloor+i==order.Floor) || (status.State=="DOWN" && status.LastFloor-i==order.Floor) && status.Inhouse==false){
//						Println("Elevator ", statuslist[elevator], "taking order: ", order)
						if statelist[elevator].Source==myip{
							return order
						}else{
							Println("BEFORE: ", networking.GetStatusList())
							Println("BEFORE: ", statelist)
							delete(connections, elevator)
//							delete(statelist,elevator)
							Println("After: ", networking.GetStatusList())
							Println("After: ", statelist)
							continue orderloop
						}
					}
				}
			}
			for elevator,_ :=range connections{
				if status,ok := statelist[elevator]; ok{
					Println("Checking:  ", statelist[elevator], "take order: ", order)
					if status.State=="IDLE" && (status.LastFloor==order.Floor+i || status.LastFloor==order.Floor-i)&& status.Inhouse==false{
//						Println("Elevator ", statuslist[elevator], "taking order: ", order)
						if statelist[elevator].Source==myip{
							return order
						}else{
							delete(connections, elevator)
//							delete(statelist,elevator)
							continue orderloop
						}
					}
				}
			}
		}
	}
	return networking.Order{Direction:0,Floor:0,InOut:0}
}


func nextstate(myip string, connections map[string]bool, mystate string) (string, networking.Order){
	next := nextorder(myip, connections)
	Println("My next order: ", next)
	if elevator.Elev_at_floor() && next.Floor==elevator.Current_floor(){
		return "DOOR_OPEN", next
	}else if next.Floor>elevator.Current_floor(){
		return "UP", next
	}else if (next.Floor<elevator.Current_floor() && next.Floor!=0){
		return "DOWN", next
	}else if elevator.Elev_at_floor(){
		return "IDLE", next
	}else{
	return mystate, networking.Order{Direction:0,Floor:0,InOut:0}}
}

func main() {

	myip := misc.GetLocalIP()
	Println(myip)

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
	var mystatus networking.Status
	var order networking.Order
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
				Println("OMGGGGG")
				networking.NewStatus(mystatus, generatedMsgs_c)
				elevator.Elev_set_speed(-300)
				state , order = nextstate(myip, connections, mystatus.State)
			}
			case "IDLE":{
				elevator.Elev_set_speed(0)
//				Println(nextorder(myip))
				state , order = nextstate(myip, connections, mystatus.State)
			}
			case "UP":{
				elevator.Elev_set_speed(300)
				state, order = nextstate(myip, connections, mystatus.State)
			}
			case "DOWN":{
				elevator.Elev_set_speed(-300)
				state, order = nextstate(myip, connections, mystatus.State)
			}
			case "DOOR_OPEN":{
				elevator.Elev_set_door_open_lamp(1)
				order.InOut=0
				networking.Neworder(generatedMsgs_c, order)
				elevator.Elev_set_speed(0)
				time.Sleep(3000 * time.Millisecond)
				elevator.Elev_set_door_open_lamp(0)
				state , order = nextstate(myip, connections, mystatus.State)
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
		time.Sleep(10 * time.Millisecond)
//		Println(state)
		elevator.FloorUpdater()
		mystatus.State=state
		mystatus.LastFloor=elevator.Current_floor()
//		mystatus.Inhouse=ConflictingOrders(mystatus, ordersinside)
		networking.NewStatus(mystatus, generatedMsgs_c)
		time.Sleep(10 * time.Millisecond)
//		generatedMsgs_c <- networking.GenerateMessage(elevator.BUTTON_CALL_UP,0,0,mystatus.State, mystatus.LastFloor,false,mystatus.Source)
	}
}
