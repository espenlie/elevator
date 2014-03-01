package networking

import (
    "fmt"
    "net"
    ."strings"
    "time"
)

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
        case msg := <-receivedMsgs_c:
            fmt.Println(msg)
        case newConn := <- newConn_c:
			go Receiver(newConn, receivedMsgs_c)
        }
    }
}

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
