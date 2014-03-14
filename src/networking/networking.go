package networking

import (
    "fmt"
    "net"
    ."strings"
    "time"
    "drivers"
)

func Orderdistr(generatedMsgs_c chan string){
    for{
        if drivers.ReadBit(drivers.FLOOR_UP1){
            generatedMsgs_c <- "Order_UP_0"+"EOL"
        }
        if drivers.ReadBit(drivers.FLOOR_UP2){
            generatedMsgs_c <- "Order_UP_1"+"EOL"
        }
        if drivers.ReadBit(drivers.FLOOR_UP3){
            generatedMsgs_c <- "Order_UP_2"+"EOL"
        }
        if drivers.ReadBit(drivers.FLOOR_DOWN2){
            generatedMsgs_c <- "Order_DOWN_1"+"EOL"
        }
        if drivers.ReadBit(drivers.FLOOR_DOWN3){
            generatedMsgs_c <- "Order_DOWN_2"+"EOL"
        }
        if drivers.ReadBit(drivers.FLOOR_DOWN4){
            generatedMsgs_c <- "Order_DOWN_3"+"EOL"
        }
	time.Sleep(5 * time.Millisecond)
    }
}

//Receives messages from a connections and adds it to a channel
func Receiver(conn net.Conn, receivedMsgs_c chan string){
    var buf [1024]byte
    for {
        _, err := conn.Read(buf[0:])
        if err != nil {
            fmt.Println(err)
            return
        }
        receivedMsgs_c <- Split(string(buf[0:]), "EOL")[0]
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


func Networking(newConn_c chan net.Conn, generatedMsgs_c chan string, receivedMsgs_c chan string, dialConn_c chan net.Conn) {
    var newConn net.Conn
    var msg []string
	connMap := make(map[string]net.Conn)
    for{
        select {
        case newConn = <- dialConn_c:
			connMap[newConn.LocalAddr().String()] = newConn
		case sendMsg := <- generatedMsgs_c:{
			for _,connection := range connMap{
				connection.Write(append([]byte(sendMsg), []byte{0}...))
			}
		}
        case in := <-receivedMsgs_c:
            fmt.Println(in)
            msg = Split(in,"_")
            if msg[0]=="Order"{
                if msg[1]=="UP" && msg[2]=="0"{
                    drivers.SetBit(drivers.LIGHT_UP1)}
                if msg[1]=="UP" && msg[2]=="1"{
                    drivers.SetBit(drivers.LIGHT_UP2)}
                if msg[1]=="UP" && msg[2]=="2"{
                    drivers.SetBit(drivers.LIGHT_UP3)}
                if msg[1]=="DOWN" && msg[2]=="1"{
                    drivers.SetBit(drivers.LIGHT_DOWN2)}
                if msg[1]=="DOWN" && msg[2]=="2"{
                    drivers.SetBit(drivers.LIGHT_DOWN3)}
                if msg[1]=="DOWN" && msg[2]=="3"{
                    drivers.SetBit(drivers.LIGHT_DOWN4)}
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
