package main
import (
    ."fmt"
    ."net"
//  ."strings"
//  "time"
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
    var sendMessage string

    go networking.Networking(newConn_c, generatedMsgs_c, receivedMsgs_c, dialConn_c)

	go networking.Listener(listenConn, newConn_c)
	go networking.Dialer(connections, conf.Default_Dial_Port, dialConn_c)

	for {
        Scanf("%s", &sendMessage)
        generatedMsgs_c <- sendMessage+"EOL" 
        }

	state := "INIT"
	for{
		Println("State: ", state)
		switch state {
    		case "INIT":{ 
				state="IDLE"
				fallthrough;
			}
    		case "IDLE":{ 
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
