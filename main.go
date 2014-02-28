package main
import (
    ."fmt"
    ."net"
    ."strings"
//  "time"
//   "drivers"
	"misc"
)

func receiver(conn Conn, receivedMsgs_c chan string){
    var buf [1024]byte
    for {
        _, err := conn.Read(buf[0:])
        if err != nil {
            Println(err)
            return
        }
        receivedMsgs_c <- Split(string(buf[0:]), "EOL")[0]
    }
}

func listener(conn *TCPListener, newConn_c chan Conn){
    for {
        newConn, err := conn.Accept()
        if err != nil {
            Println(err)
        }
        newConn_c <- newConn
    }
}

func networking(newConn_c chan Conn, generatedMsgs_c chan string, receivedMsgs_c chan string, dialConn_c chan Conn) {
    var newConn Conn
    for{
        select {
        case newConn = <- dialConn_c:
		case sendMsg := <- generatedMsgs_c:
			newConn.Write(append([]byte(sendMsg), []byte{0}...))
        case msg := <-receivedMsgs_c:
            Println(msg)
        case newConn := <- newConn_c:
//			connMap[newConn.LocalAddr().String()] = newConn
			go receiver(newConn, receivedMsgs_c)
	
        }
    }
}

func dialer(elevators map[string]bool, port string, dialconn_c chan Conn){
	for{
		for elevator,status := range elevators{
			Println(elevators)
			if !status{
				dialConn, err := Dial("tcp", elevator+port)
				if err != nil {
					Println(err)
				}else{
					elevators[elevator]=true
					dialconn_c <-dialConn
				}
			}
		}
	}
}

func main() {

//	var conf misc.Config
	conf := misc.LoadConfig("/home/student/LL/elevator/config/conf.json")
	
    connections         := make(map[string]bool)
	
	for _,elevator :=range conf.Elevators{
		Println(elevator.Address)
		connections[elevator.Address]=false
	}
	Println(connections)
    listenAddr, _ := ResolveTCPAddr("tcp", ":6969")
    listenConn, _ := ListenTCP("tcp", listenAddr)
    receivedMsgs_c  := make(chan string)
    generatedMsgs_c  := make(chan string)
    newConn_c       := make(chan Conn, 10)
    dialConn_c      := make(chan Conn, 10)
    var sendMessage string

    go networking(newConn_c, generatedMsgs_c, receivedMsgs_c, dialConn_c)

	go listener(listenConn, newConn_c)
	go dialer(connections, conf.Default_Dial_Port, dialConn_c)
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
