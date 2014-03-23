package networking

import (
	"drivers"
	"elevator"
	"encoding/json"
	"fmt"
	"log"
	"misc"
	"net"
	"os"
	"strings"
	"time"
)

type Con struct {
	Address *net.TCPConn
	Connect bool
}

var elevators = make(map[string]bool)
var orderlist = make([]Order, 0)
var statuslist = make(map[string]Status)
var connections = make([]*net.TCPConn, 0)

func GetStatusList() map[string]Status {
	return statuslist
}

func GetOrderList() []Order {
	return orderlist
}

func PackNetworkMessage(message Networkmessage, error_c chan string) []byte {
	send, err := json.Marshal(message)
	if err != nil {
		error_c <- "Could not pack message: " + err.Error()
	}
	return send
}

func UnpackNetworkMessage(pack []byte, error_c chan string) Networkmessage {
	var message Networkmessage
	err := json.Unmarshal(pack, &message)
	if err != nil {
		error_c <- "Could not unpack message: " + err.Error()
	}
	return message
}

//Generates and sends an initialization message to other elevators
func InitUpdate(connection *net.TCPConn, myip string, error_c chan string) {
	pack := make([]byte, 1024)
	status := statuslist[myip]
	pack = PackNetworkMessage(Networkmessage{Order: Order{}, Status: status}, error_c)
	time.Sleep(10 * time.Millisecond)
	connection.Write(pack)
	for _, order := range orderlist {
		time.Sleep(10 * time.Millisecond)
		pack = PackNetworkMessage(Networkmessage{Order: order, Status: Status{}}, error_c)
		connection.Write(pack)
	}
}

//Check if there is any new orders, if it is it passes it to Neworder
func Orderdistr(generatedMsgs_c chan Networkmessage, myip string) {
	var butt elevator.Elev_button
	for {
		for floor, buttons := range elevator.Button_channel_matrix {
			for button, channel := range buttons {
				if drivers.ReadBit(channel) {
					if button == 0 {
						butt = elevator.BUTTON_CALL_UP
					} else if button == 1 {
						butt = elevator.BUTTON_CALL_DOWN
					} else {
						butt = elevator.BUTTON_COMMAND
					}
					Neworder(generatedMsgs_c, Order{Direction: butt, Floor: floor + 1, InOut: 1, Source: myip})
					time.Sleep(time.Millisecond)
				}
			}
		}
	}
}

//Dials elevators in config that we do not have connection with
func Dialer(connect_c chan Con, port string, elevators []misc.Elevator, error_c chan string) {
	local, _ := net.ResolveTCPAddr("tcp", "localhost"+port)
	localconn, _ := net.DialTCP("tcp", nil, local)
	connect_c <- Con{Address: localconn, Connect: true}
	for {
	elevatorloop:
		for _, elevator := range elevators {
			cons := connections
			for _, connection := range cons {
				if strings.Split(connection.RemoteAddr().String(), ":")[0] == elevator.Address {
					continue elevatorloop
				}
			}
			raddr, err := net.ResolveTCPAddr("tcp", elevator.Address+port)
			dialConn, err := net.DialTCP("tcp", nil, raddr)
			if err != nil {
				error_c <- "Dial trouble: " + err.Error()
			} else {
				connect_c <- Con{Address: dialConn, Connect: true}
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func Listener(conn *net.TCPListener, connect_c chan Con, error_c chan string) {
	for {
		newConn, err := conn.AcceptTCP()
		if err != nil {
			error_c <- "Accept trouble: " + err.Error()
		}
		connect_c <- Con{Address: newConn, Connect: true}
	}
}

func Receiver(conn *net.TCPConn, receivedMsgs_c chan Networkmessage, connections_c chan Con, error_c chan string) {
	buf := make([]byte, 1024)
	keepalivebyte := []byte("KEEPALIVE")
receiverloop:
	for {
		err := conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		if err != nil {
			error_c <- "Trouble setting read deadline: " + err.Error()
			connections_c <- Con{Address: conn, Connect: false}
			return
		}
		bit, err := conn.Read(buf[0:])
		if err != nil {
			error_c <- "Trouble receiving: " + err.Error()
			connections_c <- Con{Address: conn, Connect: false}
			return
		}
		if string(buf[:bit]) == string(keepalivebyte) {
			continue receiverloop
		}
		unpacked := UnpackNetworkMessage(buf[:bit], error_c)
		receivedMsgs_c <- unpacked
	}
}

func SendAliveMessages(connection *net.TCPConn, error_c chan string) {
	for {
		_, err := connection.Write([]byte("KEEPALIVE"))
		if err != nil {
			error_c <- "Problems sending keepalive message: " + err.Error()
			return
		}
		time.Sleep(time.Second)
	}
}

//Function that controls all the network communication and connections
func TCPPeerToPeer(conf misc.Config, myip string, generatedmessages_c chan Networkmessage) {
	elevlog, err := os.OpenFile("elevator.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file: " + err.Error())
	}
	defer elevlog.Close()
	log.SetOutput(elevlog)
	listenaddr, _ := net.ResolveTCPAddr("tcp", conf.DefaultListenPort)
	listenconn, _ := net.ListenTCP("tcp", listenaddr)
	connections_c := make(chan Con, 15)
	receivedmessages_c := make(chan Networkmessage, 15)
	error_c := make(chan string, 10)
	go Listener(listenconn, connections_c, error_c)
	go Orderdistr(generatedmessages_c, myip)
	go Dialer(connections_c, conf.DefaultListenPort, conf.Elevators, error_c)
	for {
		select {
		case connection := <-connections_c: //Managing new/closed connections
			{
				if connection.Connect {
					connections = append(connections, connection.Address)
					go Receiver(connection.Address, receivedmessages_c, connections_c, error_c)
					go SendAliveMessages(connection.Address, error_c)
					go InitUpdate(connection.Address, myip, error_c)
				} else {
					remoteip := strings.Split(connection.Address.RemoteAddr().String(), ":")[0]
					errorstate := Status{State: "ERROR", LastFloor: 0, Inhouse: false, Source: remoteip}
					statuslist[remoteip] = errorstate
					for i, con := range connections {
						if con == connection.Address {
							connections[len(connections)-1], connections[i], connections = nil, connections[len(connections)-1], connections[:len(connections)-1]
						}
					}
					connection.Address.Close()
				}

			}
		case received := <-receivedmessages_c:
			{
				if received.Order.Floor > 0 {
					if !((received.Order.Direction == elevator.BUTTON_COMMAND) && (received.Order.Source != myip)) {
						elevator.ElevSetButtonLamp(received.Order.Direction, received.Order.Floor, received.Order.InOut)
					}
					if received.Order.Direction != elevator.BUTTON_COMMAND {
						received.Order.Source = ""
					}
					if received.Order.InOut == 0 {
						received.Order.InOut = 1
						for i, b := range orderlist {
							if b == received.Order {
								orderlist = append(orderlist[:i], orderlist[i+1:]...)
							}
						}
					} else {
						AddedBefore := false
						for _, b := range orderlist {
							if b == received.Order {
								AddedBefore = true
							}
						}
						if !AddedBefore {
							orderlist = append(orderlist, received.Order)
						}
					}
				}
				if received.Status.Source != "" {
					statuslist[received.Status.Source] = received.Status
				}
			}
		case message := <-generatedmessages_c:
			{
				pack := make([]byte, 1024)
				pack = PackNetworkMessage(message, error_c)
				for _, connection := range connections {
					_, err := connection.Write(pack)
					if err != nil {
						error_c <- "Problems writing to connection: " + err.Error()
					}
				}
			}
		case err := <-error_c:
			{
				log.Println("ERROR: " + err)
			}
		}
	}
}

func SendStatuslist(generatedMsgs_c chan Networkmessage) {
	myip := misc.GetLocalIP()
	mystatus := statuslist[myip]
	generatedMsgs_c <- Networkmessage{Order: Order{}, Status: mystatus}
}

func NewStatus(status Status, generatedMsgs_c chan Networkmessage) bool {
	for _, oldstat := range statuslist {
		if oldstat == status {
			return false
		}
	}
	generatedMsgs_c <- Networkmessage{Order: Order{}, Status: status}
	return true
}

func Neworder(generatedMsgs_c chan Networkmessage, order Order) bool {
	if order.Direction != elevator.BUTTON_COMMAND {
		order.Source = ""
	}
	for _, b := range orderlist {
		if b == order {
			return false
		}
	}
	generatedMsgs_c <- Networkmessage{Order: order, Status: Status{}}
	return true
}
