package main
import (
    ."fmt"
    ."net"
//  ."strings"
//	"strconv"
    "time"
    "drivers"
	"misc"
	"elevator"
	"networking"
)

var orderlist = make([]networking.Order,0)
var statuslist = make(map[string]networking.Status)

func takeorder(orderlist []networking.Order, statuslist map[string]networking.Status, myip string)int{
	Println(orderlist, statuslist, myip)
//	for order := orderlist.Front(); order != nil; order = order.Next() {
//		for elevator := statuslist.Front(); elevator != nil; elevator = elevator.Next() {
//			if elevator.State=="IDLE" && elevator.LastFloor==order.Floor && elevator.Source=myip{
//				return order.Floor
//			}
//			if elevator.State=="UP" && elevator.LastFloor==order.Floor-1 && elevator.Source=myip{
//				return order.Floor
//			}
//		}
//	}
	return -1
}

func main() {

	var myip string
	myip = "yooo"
	drivers.IoInit()
	elevator.Elev_init()
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
	takeorder(orderlist, statuslist, myip)
	go networking.Listener(listenConn, newConn_c)
	go networking.Dialer(connections, conf.Default_Dial_Port, dialConn_c)
	go networking.Orderdistr(generatedMsgs_c)
//	for {
//      Scanf("%s", &sendMessage)
//      generatedMsgs_c <- sendMessage+"EOL" 

//      }

	state := "INIT"
	var floor int
	var status networking.Status
	interfaces, _ := Interfaces()
		for _, inter := range interfaces {
			if inter.Name=="rename2"{
				addrs, _ := inter.Addrs()
				for _, addr := range addrs {
					status.Source=addr.String()
				}
			}
		}
	Println(status.Source)

	for{
//		Println("State: ", state)
		switch state {
			case "INIT":{
				state="IDLE"
				fallthrough;
				}
			case "IDLE":{ 
				if drivers.ReadBit(drivers.SENSOR1){
					floor=1
				}
				if drivers.ReadBit(drivers.SENSOR2){
					floor=2
				}
				if drivers.ReadBit(drivers.SENSOR3){
					floor=3
				}
				if drivers.ReadBit(drivers.SENSOR4){
					floor=4
				}
				time.Sleep(100 * time.Millisecond)
				status.State = state
				status.LastFloor = floor
//      		generatedMsgs_c <- status.Source+"_"+status.State+"_"+strconv.Itoa(status.LastFloor)+"EOL"
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
			case "ESTOP":{ 
				fallthrough;
			}
			case "ERROR":{ 
				fallthrough;
			}
			default:{

//  			generatedMsgs_c <- "Jeg lever"
			}
    	}
	}
}
