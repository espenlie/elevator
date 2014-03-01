package main
import (
    ."fmt"
    ."net"
//  ."strings"
	"strconv"
    "time"
    "drivers"
	"misc"
	"elevator"
	"networking"
)

func main() {
	drivers.IoInit()
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
    receivedMsgs_c  := make(chan string)
    generatedMsgs_c  := make(chan string)
    newConn_c       := make(chan Conn, 10)
    dialConn_c      := make(chan Conn, 10)

    go networking.Networking(newConn_c, generatedMsgs_c, receivedMsgs_c, dialConn_c)
	
	go networking.Listener(listenConn, newConn_c)
	go networking.Dialer(connections, conf.Default_Dial_Port, dialConn_c)

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
        		generatedMsgs_c <- status.Source+"_"+status.State+"_"+strconv.Itoa(status.LastFloor)+"EOL"
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
