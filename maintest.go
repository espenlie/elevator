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


func Nextorder(myip string, Elevatorlist []Elevator)networking.Order{
	var statelist = make(map[string]networking.Status)
	statuslist := networking.GetStatusList()
	insidelist := networking.GetInsidelist()
	for host, status := range statuslist {
		statelist[host]=status
    }
	orderlist := networking.GetOrderList()
//	Println("orderlist: ", orderlist)
//	Println("Statelist: ", statelist)
//	Println("Connections: ", conf.Elevators)
//	Println("MYIP: ", myip)
	insideloop:
	for _,order := range insidelist{
		for _, elevator :=range Elevatorlist{
			if status,ok := statelist[elevator.Address]; ok{
				if ((status.State=="UP" && status.LastFloor<=order.Floor) || (status.State=="DOWN" && status.LastFloor>=order.Floor)){
					if statelist[elevator.Address].Source==myip{
						return order
					}else{
						delete(statelist,elevator.Address)
						continue insideloop
					}
				}
			}
		}
		for _, elevator :=range Elevatorlist{
			if status,ok := statelist[elevator.Address]; ok{
				if ((status.State=="UP" && status.LastFloor>=order.Floor) || (status.State=="DOWN" && status.LastFloor<=order.Floor)){
					if statelist[elevator.Address].Source==myip{
						return order
					}else{
						delete(statelist,elevator.address)
						continue insideloop
					}
				}
			}
		}
	}
	orderloop:
	for _,order := range orderlist{
		for i := 0; i < elevator.N_FLOORS; i++ {
			for _, elevator :=range Elevatorlist{
				if status,ok := statelist[elevator.Address]; ok{
					if (i!=0 && (status.State=="UP" && status.LastFloor+i==order.Floor) || (status.State=="DOWN" && status.LastFloor-i==order.Floor)){
						if statelist[elevator.Address].Source==myip{
							return order
						}else{
							delete(statelist,elevator.address)
							continue orderloop
						}
					}
				}
			}
			for _, elevator :=range Elevatorlist{
				if status,ok := statelist[elevator.Address]; ok{
					if status.State=="IDLE" && (status.LastFloor==order.Floor+i || status.LastFloor==order.Floor-i){
						if statelist[elevator.Address].Source==myip{
							return order
						}else{
							delete(statelist,elevator.address)
							continue orderloop
						}
					}
				}
			}
		}
	}
	return networking.EmptyOrder
}

func Stop(myip string, mystate string)networking.Order{
	takeorder := networking.EmptyOrder
	insidelist := networking.GetInsidelist()
	orderlist := networking.GetOrderList()
	for _,order := range insidelist{
		if order.Source==myip && order.Floor==elevator.Current_floor(){
			takeorder=append(takeorder, order)
		}
	}
	for _,order := range orderlist{
		if ((order.Direction==elevator.BUTTON_CALL_UP && mystate=="UP") || (order.Direction==elevator.BUTTON_CALL_DOWN && mystate=="DOWN")){
			if order.Floor==elevator.Current_floor(){
				takeorder=append(takeorder, order)
			}
		}
	}
	return takeorder
}


func nextstate(myip string, conf.Elevators []misc.Elevator, mystate string) (string, []networking.Order){
	Println("My next order: ", next)
	stop := Stop(myip, mystate)
	if elevator.Elev_at_floor() && stop.Floor==elevator.Current_floor(){
		return "DOOR_OPEN", networking.Order{Direction:0,Floor:0,InOut:0}
//	}else if elevator.Elev_at_floor() && next.Floor==elevator.Current_floor(){  //Behoves denne?
//		Println("DENNE KAN IKKE SLETTES!")
//		return "DOOR_OPEN", next												//Behoves denne?
	next := Nextorder(myip, conf.Elevators)
	}else if next.Floor>elevator.Current_floor(){
		return "UP", next
	}else if (next.Floor<elevator.Current_floor() && next.Floor!=0){
		return "DOWN", next
	}else if elevator.Elev_at_floor(){
		return "IDLE", next
	}else{
	return mystate, []networking.EmptyOrder
	}
}

func main() {

	myip := misc.GetLocalIP()
	Println(myip)

//	var conf misc.Config
	conf := misc.LoadConfig("/home/student/LL/elevator/config/conf.json")

    connections         := make(map[string]bool)


//  listenAddr, _ := ResolveTCPAddr("tcp", ":6969")
//  listenConn, _ := ListenTCP("tcp", listenAddr)
//  receivedMsgs_c  := make(chan networking.Networkmessage, 10)
//  generatedMsgs_c  := make(chan networking.Networkmessage, 10)
//  newConn_c       := make(chan *TCPConn, 10)
//  dialConn_c      := make(chan *TCPConn, 10)

    go networking.Networking(conf)
//	statuslist[myip]=networking.Status{State:"UP",LastFloor:1,Source:myip}
//	takeorder(orderlist, statuslist, myip)
//	go networking.Listener(listenConn, newConn_c)
//	go networking.Dialer(connections, conf.Default_Dial_Port, dialConn_c)
//	go networking.Orderdistr(generatedMsgs_c)
//	for {
//      Scanf("%s", &sendMessage)
//      generatedMsgs_c <- sendMessage+"EOL" 

//      }

	state := "INIT"
//	var floor int
	var mystatus networking.Status
	var takeorders []networking.Order
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
				state , takeorders = nextstate(myip, conf.Elevators, mystatus.State)
			}
			case "IDLE":{
				elevator.Elev_set_speed(0)
				state , takeorders = nextstate(myip, conf.Elevators, mystatus.State)
			}
			case "UP":{
				elevator.Elev_set_speed(300)
				state, takeorders = nextstate(myip, conf.Elevators, mystatus.State)
			}
			case "DOWN":{
				elevator.Elev_set_speed(-300)
				state, takeorders = nextstate(myip, conf.Elevators, mystatus.State)
			}
			case "DOOR_OPEN":{
				elevator.Elev_set_door_open_lamp(1)
				for _, order = range takeorders{
					order.InOut=0
					networking.Neworder(generatedMsgs_c, order)
				}
				elevator.Elev_set_speed(0)
				time.Sleep(3000 * time.Millisecond)
				elevator.Elev_set_door_open_lamp(0)
				state , order = nextstate(myip, conf.Elevators, mystatus.State)
			}
			case "ERROR":{
			}
		}    
//		Println(elevator.Address)
//		statuslist := networking.GetStatusList()
//		orderlist := networking.GetOrderList()
//		Println(statuslist)
//		Println(orderlist)
//		Println(state)
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
